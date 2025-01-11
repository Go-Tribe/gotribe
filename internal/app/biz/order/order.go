package order

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/log"
	util "gotribe/pkg/amount"
	"gotribe/pkg/api/v1"

	"github.com/jinzhu/copier"
	"gotribe/internal/app/store"
	"gotribe/internal/pkg/model"
)

// OrderBiz defines functions used to handle comment request.
type OrderBiz interface {
	CreateTx(ctx context.Context, username string, r *v1.CreateOrderRequest) (string, error)
	Get(ctx context.Context, orderNumber, username string) (*v1.GetOrderResponse, error)
	List(ctx context.Context, username string, offset, limit int) (*v1.ListOrderResponse, error)
}

// The implementation of OrderBiz interface.
type orderBiz struct {
	ds store.IStore
}

// Make sure that orderBiz implements the OrderBiz interface.
// We can find this problem in the compile stage with the following assignment statement.
var _ OrderBiz = (*orderBiz)(nil)

func New(ds store.IStore) *orderBiz {
	return &orderBiz{ds: ds}
}

// CreateTx 实现事务的示例
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

func (b *orderBiz) Get(ctx context.Context, orderNumber, username string) (*v1.GetOrderResponse, error) {
	order, err := b.ds.Order().Get(ctx, v1.OrderWhere{OrderNumber: orderNumber, Username: username})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrOrderNotFound
		}

		return nil, err
	}

	var resp v1.GetOrderResponse
	_ = copier.Copy(&resp, order)
	resp.Amount = util.FenToYuan(order.Amount)
	resp.AmountPay = util.FenToYuan(order.AmountPay)
	resp.UnitPoint = util.FenToYuan(order.UnitPoint)
	resp.CreatedAt = order.CreatedAt.Format(known.TimeFormat)
	resp.UpdatedAt = order.UpdatedAt.Format(known.TimeFormat)
	resp.PayTime = order.PayTime.Format(known.TimeFormat)

	return &resp, nil
}

func (b *orderBiz) List(ctx context.Context, username string, offset, limit int) (*v1.ListOrderResponse, error) {
	count, list, err := b.ds.Order().List(ctx, offset, limit, v1.OrderWhere{Username: username, ProjectID: ctx.Value(known.XPrjectIDKey).(string)})
	if err != nil {
		log.C(ctx).Errorw("Failed to list order from storage", "err", err)
		return nil, err
	}

	orders := make([]*v1.OrderInfo, 0, len(list))
	for _, item := range list {
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

	return &v1.ListOrderResponse{TotalCount: count, Orders: orders}, nil
}
