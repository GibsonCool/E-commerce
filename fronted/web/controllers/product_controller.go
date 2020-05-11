package controllers

import (
	"E-commerce/services"
	"github.com/kataras/iris"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
}
