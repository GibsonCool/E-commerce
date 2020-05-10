package controllers

import (
	"E-commerce/common/conf"
	"E-commerce/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

type ProductController struct {
	Ctx          iris.Context
	OrderService services.IProductService
}

func (pc *ProductController) GetAll() mvc.View {
	product, err := pc.OrderService.GetAllProduct()
	if err != nil {
		conf.AppSetting.Logger.Error(err.Error())
	}

	return mvc.View{
		Name: "product/view.html",
		Data: iris.Map{
			"productArray": product,
		},
	}
}
