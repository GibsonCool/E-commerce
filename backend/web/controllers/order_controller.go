package controllers

import (
	"E-commerce/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/unknwon/com"
)

type OrderController struct {
	Ctx          iris.Context
	OrderService services.IOrderService
}

func (oc *OrderController) Get() mvc.View {
	info, err := oc.OrderService.GetAllOrderInfo()
	if err != nil {
		oc.Ctx.Application().Logger().Error("查询订单失败：" + err.Error())
	}

	result := map[int]map[string]string{}
	for index, orderInfo := range info {
		result[index] = map[string]string{
			"ID":          com.ToStr(orderInfo.ID),
			"userName":    orderInfo.UserName,
			"productName": orderInfo.ProductName,
			"orderStatus": orderInfo.OrderStatus,
		}
	}

	return mvc.View{
		Name: "order/view.html",
		Data: iris.Map{
			"order": result,
		},
	}
}
