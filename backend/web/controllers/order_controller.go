package controllers

import (
	"E-commerce/services"
	"github.com/kataras/iris"
)

type OrderController struct {
	Ctx          iris.Context
	OrderService services.IOrderService
}

//func (oc *OrderController) Get() mvc.View {
//
//}
