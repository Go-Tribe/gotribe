package order

import (
	"context"
	"github.com/dengmengmian/ghelper/gid"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gotribe/internal/app/biz/point"
	"gotribe/internal/app/store"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/log"
	"gotribe/internal/pkg/model"
	util "gotribe/pkg/amount"
	"gotribe/pkg/api/v1"
	"time"
)

// OrderBiz defines functions used to handle comment request.
type OrderBiz interface {
	CreateTx(ctx context.Context, username string, r *v1.CreateOrderRequest) (string, error)
	Get(ctx context.Context, orderNumber, username string) (*v1.GetOrderResponse, error)
	List(ctx context.Context, username string, offset, limit int) (*v1.ListOrderResponse, error)
	Pay(ctx context.Context, orderNumber, username string) error
}

// The implementation of OrderBiz interface.
type orderBiz struct {
	ds store.IStore
	pt point.PointBiz
}

// Make sure that orderBiz implements the OrderBiz interface.
// We can find this problem in the compile stage with the following assignment statement.
var _ OrderBiz = (*orderBiz)(nil)

func New(ds store.IStore, pt point.PointBiz) *orderBiz {
	return &orderBiz{ds: ds, pt: pt}
}

// CreateTx 创建订单并处理事务
// 该方法在事务中执行多个步骤，确保数据的一致性
// 参数:
//
//	ctx - 上下文，用于传递请求范围的数据
//	username - 用户名，用于获取用户信息
//	r - 创建订单的请求，包含订单的相关信息
//
// 返回值:
//
//	string - 创建的订单号
//	error - 错误信息，如果执行过程中发生错误
func (b *orderBiz) CreateTx(ctx context.Context, username string, r *v1.CreateOrderRequest) (string, error) {
	var orderNumber string
	err := b.ds.TX(ctx, func(ctx context.Context) error {
		// 1. 获取用户信息
		user, err := b.ds.Users().Get(ctx, v1.UserWhere{Username: username})
		if err != nil {
			return err
		}
		// 2. 获取商品信息
		product, err := b.ds.Products().Get(ctx, r.ProductID)
		if err != nil {
			return err
		}
		// 3. 获取商品SKU信息
		sku, err := b.ds.ProductSKUs().Get(ctx, r.ProductSkuID)
		if err != nil {
			return err
		}
		if sku.Quantity < uint(r.Quantity) {
			return errors.New("库存不够")
		}

		var orderM model.OrderM
		if err := copier.Copy(&orderM, r); err != nil {
			return err
		}
		orderM.UserID = user.UserID
		orderM.Username = username
		orderM.ProductSku = sku.SkuID
		orderM.ProductName = product.Title
		orderM.Status = known.OrderStatusPendingPayment
		orderM.ProjectID = ctx.Value(known.XPrjectIDKey).(string)
		orderM.UnitPoint = sku.UnitPoint
		orderM.UnitPrice = sku.UnitPrice
		orderM.Quantity = uint(r.Quantity)
		orderM.Amount = sku.UnitPoint * r.Quantity
		orderM.AmountPay = sku.UnitPrice * r.Quantity

		createdOrder, err := b.ds.Order().Create(ctx, &orderM)
		if err != nil {
			return err
		}
		// 扣减库存
		if err := b.ds.ProductSKUs().Update(ctx, &model.ProductSKUM{
			Quantity: sku.Quantity - uint(r.Quantity),
		}); err != nil {
			return err
		}
		orderNumber = createdOrder.OrderNumber
		return nil
	})
	if err != nil {
		return "", err
	}
	return orderNumber, nil
}

// Get 通过订单号和用户名获取订单详情。
// 该方法首先调用数据源（ds）的 Order().Get() 方法来查询数据库中的订单信息。
// 如果查询过程中出现错误，特别是当记录不存在时，方法会返回一个标准的错误号（ErrOrderNotFound）。
// 查询成功后，将数据库中的订单信息转换为 GetOrderResponse 对象，并对某些字段进行格式化处理，
// 如将金额从分转换为元，将日期时间格式化为已知的时间格式，最后返回该响应对象。
func (b *orderBiz) Get(ctx context.Context, orderNumber, username string) (*v1.GetOrderResponse, error) {
	// 从数据源获取特定订单号和用户名对应的订单信息。
	order, err := b.ds.Order().Get(ctx, v1.OrderWhere{OrderNumber: orderNumber, Username: username})
	if err != nil {
		// 当订单记录不存在时，返回自定义的错误号。
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrOrderNotFound
		}

		// 返回查询过程中遇到的其他错误。
		return nil, err
	}

	// 初始化 GetOrderResponse 对象，并将查询到的订单信息复制到该对象中。
	var resp v1.GetOrderResponse
	_ = copier.Copy(&resp, order)

	// 将订单金额从分转换为元，并更新到响应对象中。
	resp.Amount = util.FenToYuan(order.Amount)
	resp.AmountPay = util.FenToYuan(order.AmountPay)
	resp.UnitPoint = util.FenToYuan(order.UnitPoint)

	// 将订单的创建时间、更新时间和支付时间格式化为已知的时间格式。
	resp.CreatedAt = order.CreatedAt.Format(known.TimeFormat)
	resp.UpdatedAt = order.UpdatedAt.Format(known.TimeFormat)
	resp.PayTime = order.PayTime.Format(known.TimeFormat)

	// 返回填充好的 GetOrderResponse 对象。
	return &resp, nil
}

// List 获取订单列表
// 该方法根据用户名、偏移量和限制数量查询订单，并返回订单列表和总订单数
// 主要用于处理订单数据的获取，包括从数据源获取原始数据，然后将这些数据转换为API响应格式
func (b *orderBiz) List(ctx context.Context, username string, offset, limit int) (*v1.ListOrderResponse, error) {
	// 调用数据源的List方法获取订单数据
	count, list, err := b.ds.Order().List(ctx, offset, limit, v1.OrderWhere{Username: username, ProjectID: ctx.Value(known.XPrjectIDKey).(string)})
	if err != nil {
		// 如果查询失败，记录错误日志并返回错误
		log.C(ctx).Errorw("Failed to list order from storage", "err", err)
		return nil, err
	}

	// 初始化订单信息切片
	orders := make([]*v1.OrderInfo, 0, len(list))
	for _, item := range list {
		// 遍历查询结果，将每个订单转换为API响应格式
		order := item
		orders = append(orders, &v1.OrderInfo{
			Username:          order.Username,
			OrderID:           order.OrderID,
			OrderNumber:       order.OrderNumber,
			PayNumber:         order.PayNumber,
			PayTime:           order.PayTime.Format(known.TimeFormat),
			Status:            uint8(order.Status),
			Amount:            util.FenToYuan(order.Amount),
			AmountPay:         util.FenToYuan(order.AmountPay),
			ProductID:         order.ProductID,
			ProductName:       order.ProductName,
			Quantity:          int(order.Quantity),
			UnitPoint:         util.FenToYuan(order.UnitPoint),
			ConsigneeName:     order.ConsigneeName,
			ConsigneePhone:    order.ConsigneePhone,
			ConsigneeAddress:  order.ConsigneeAddress,
			ConsigneeProvince: order.ConsigneeProvince,
			ConsigneeCity:     order.ConsigneeCity,
			ConsigneeDistrict: order.ConsigneeDistrict,
			ConsigneeStreet:   order.ConsigneeStreet,
			Remark:            order.Remark,
			CreatedAt:         order.CreatedAt.Format(known.TimeFormat),
			UpdatedAt:         order.UpdatedAt.Format(known.TimeFormat),
		})
	}

	// 返回订单列表和总订单数
	return &v1.ListOrderResponse{TotalCount: count, Orders: orders}, nil
}

// Pay 订单支付业务处理函数
// 该函数通过上下文context、订单号orderNumber和用户名username作为参数
// 它尝试执行一个事务，该事务包括以下几个步骤：
// 1. 获取当前订单信息
// 2. 获取用户信息并检查用户积分是否足够支付订单金额
// 3. 扣减用户积分并更新订单状态为已支付
// 如果任何步骤失败，事务将回滚并返回错误
func (b *orderBiz) Pay(ctx context.Context, orderNumber, username string) error {
	err := b.ds.TX(ctx, func(ctx context.Context) error {
		// 1.获取当前订单信息
		order, err := b.ds.Order().Get(ctx, v1.OrderWhere{OrderNumber: orderNumber, Username: username})
		if err != nil {
			// 如果查询失败，记录错误日志并返回错误
			log.C(ctx).Errorw("Failed to get order from storage", "err", err)
			return err
		}
		// 2. 获取用户信息
		user, err := b.ds.Users().Get(ctx, v1.UserWhere{Username: username})
		if err != nil {
			return err
		}
		// 检查用户积分是否足够支付订单金额
		if user.Point < float64(order.AmountPay) {
			return errors.New("积分不足")
		}
		// 3.扣减用户积分
		err = b.pt.SubPoints(ctx, username, order.ProjectID, "pay", "支付订单", order.OrderID, util.FenToYuan(order.AmountPay))
		if err != nil {
			return err
		}
		// 更新订单状态为已支付
		order.Status = known.OrderStatusPaid
		order.PayMethod = 3
		order.PayTime = time.Now()
		order.PayNumber = gid.FetchOrderNum(6)
		log.C(ctx).Infow("pay order", "order", order)
		// 更新数据库中的订单信息
		if err := b.ds.Order().Update(ctx, order); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
