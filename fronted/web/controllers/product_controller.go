package controllers

import (
	"E-commerce/datamodels"
	"E-commerce/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/unknwon/com"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
	OrderService   services.IOrderService
}

// 获取订单详情页
// /product/detail?id=xxx
// get
func (pc *ProductController) GetDetail() mvc.View {
	id := pc.Ctx.URLParam("id")
	pc.Ctx.Application().Logger().Error("id:" + id)
	product, err := pc.ProductService.GetProductByID(com.StrTo(id).MustInt64())
	if err != nil {
		return mvc.View{
			Name: "shared/error.html",
			Data: iris.Map{
				"Message": "获取商品信息出错：" + err.Error(),
			},
		}
	}

	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/view.html",
		Data: iris.Map{
			"product": product,
		},
	}

}

// 下单抢购
// /product/order?productID=xxx
// get
func (pc *ProductController) GetOrder() mvc.View {
	productStr := pc.Ctx.URLParam("productID")
	userIdStr := pc.Ctx.GetCookie("uid")

	// 查询商品信息
	product, err := pc.ProductService.GetProductByID(com.StrTo(productStr).MustInt64())
	if err != nil {
		return mvc.View{
			Name: "shared/error.html",
			Data: iris.Map{
				"Message": "下单抢购时，获取商品信息出错：" + err.Error(),
			},
		}
	}

	var orderId int64
	showMsg := ""
	// 判断商品数量是否足够
	if product.ProductNum > 0 {

		// 这里涉及两个表的操作，正常应该使用一个事物来执行，我们简化处理就没有使用
		// 1.商品数量 -1
		product.ProductNum -= 1
		if err := pc.ProductService.UpdateProduct(product); err != nil {
			showMsg = "抢购失败，商品查询失败：" + err.Error()
		}

		// 2.创建订单
		orderId, err = pc.OrderService.InsertOrder(&datamodels.Order{
			UserId:      com.StrTo(userIdStr).MustInt64(),
			ProductId:   int64(product.ID),
			OrderStatus: datamodels.OrderSuccess,
		})

		if err != nil {
			showMsg = "抢购失败，下单失败：" + err.Error()
		} else {
			showMsg = "抢购成功，订单号：" + com.ToStr(orderId)
		}
	}

	return mvc.View{
		Name:   "product/result.html",
		Layout: "shared/productLayout.html",
		Data: iris.Map{
			"orderID":     orderId,
			"showMessage": showMsg,
		},
	}
}
