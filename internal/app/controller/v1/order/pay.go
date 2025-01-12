package order

import (
	"gotribe/internal/pkg/core"
	"gotribe/internal/pkg/errno"
	"gotribe/internal/pkg/known"
	"gotribe/internal/pkg/log"
	"gotribe/pkg/api/v1"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// Create 创建订单
func (ctrl *OrderController) Pay(c *gin.Context) {
	log.C(c).Infow("pay order function called")

	var r v1.PayOrderRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}

	if _, err := govalidator.ValidateStruct(r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)

		return
	}

	err := ctrl.b.Orders().Pay(c, r.OrderNumber, c.GetString(known.XUsernameKey))
	if err != nil {
		core.WriteResponse(c, err, nil)

		return
	}

	core.WriteResponse(c, nil, nil)
}
