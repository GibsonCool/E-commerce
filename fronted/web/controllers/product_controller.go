package controllers

import (
	"E-commerce/datamodels"
	"E-commerce/rabbitmq"
	"E-commerce/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/unknwon/com"
	"os"
	"path/filepath"
	"text/template"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
	OrderService   services.IOrderService
	RabbitMq       *_0_RabbitMQ.RabbitMQ
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
//func (pc *ProductController) GetOrder() mvc.View {
//	productStr := pc.Ctx.URLParam("productID")
//	userIdStr := pc.Ctx.GetCookie("uid")
//
//	// 查询商品信息
//	product, err := pc.ProductService.GetProductByID(com.StrTo(productStr).MustInt64())
//	if err != nil {
//		return mvc.View{
//			Name: "shared/error.html",
//			Data: iris.Map{
//				"Message": "下单抢购时，获取商品信息出错：" + err.Error(),
//			},
//		}
//	}
//
//	var orderId int64
//	showMsg := ""
//	// 判断商品数量是否足够
//	if product.ProductNum > 0 {
//
//		// 这里涉及两个表的操作，正常应该使用一个事物来执行，我们简化处理就没有使用
//		// 1.商品数量 -1
//		product.ProductNum -= 1
//		if err := pc.ProductService.UpdateProduct(product); err != nil {
//			showMsg = "抢购失败，商品查询失败：" + err.Error()
//		}
//
//		// 2.创建订单
//		orderId, err = pc.OrderService.InsertOrder(&datamodels.Order{
//			UserId:      com.StrTo(userIdStr).MustInt64(),
//			ProductId:   int64(product.ID),
//			OrderStatus: datamodels.OrderSuccess,
//		})
//
//		if err != nil {
//			showMsg = "抢购失败，下单失败：" + err.Error()
//		} else {
//			showMsg = "抢购成功，订单号：" + com.ToStr(orderId)
//		}
//	}
//
//	return mvc.View{
//		Name:   "product/result.html",
//		Layout: "shared/productLayout.html",
//		Data: iris.Map{
//			"orderID":     orderId,
//			"showMessage": showMsg,
//		},
//	}
//}

// 将下单信息投递到 rabbitmq 让消费端执行下单操作
func (pc *ProductController) GetOrder() []byte {
	productStr := pc.Ctx.URLParam("productID")
	userIdStr := pc.Ctx.GetCookie("uid")

	message := datamodels.NewMessage(com.StrTo(productStr).MustInt64(), com.StrTo(userIdStr).MustInt64())

	if err := pc.RabbitMq.PublishSimple(message.JsonToStr()); err != nil {
		return []byte("false 加入消息队列失败：" + err.Error())
	}

	return []byte("true")
}

var (
	//生成的Html保存目录
	htmlOutPath = "./fronted/web/htmlProductShow/"
	//静态文件模版目录
	templatePath = "./fronted/web/views/template/"
)

func (pc *ProductController) GetGenerateHtml() {
	idStr := pc.Ctx.URLParam("id")
	id := com.StrTo(idStr).MustInt64()

	getWd, _ := os.Getwd()
	pc.Ctx.Application().Logger().Error("当前路径：" + getWd)

	// 1.获取模板文件
	contentTmp, err := template.ParseFiles(filepath.Join(getWd, templatePath, "product.html"))
	if err != nil {
		pc.Ctx.Application().Logger().Error("获取模板文件失败：" + err.Error())
	}
	// 2.获取模板生成路径
	fileName := filepath.Join(getWd, htmlOutPath, "htmlProduct.html")

	// 业务查询数据
	product, err := pc.ProductService.GetProductByID(id)
	if err != nil {
		pc.Ctx.Application().Logger().Error("数据获取失败：" + err.Error())
	}
	generateStaticHtml(pc.Ctx, contentTmp, fileName, product)
}

func generateStaticHtml(ctx iris.Context, template *template.Template, fileName string, product *datamodels.Product) {
	// 1.判断文件是否已经存在
	if exist(fileName) {
		// 存在则删除
		if err := os.Remove(fileName); err != nil {
			ctx.Application().Logger().Error("文件删除失败：" + err.Error())
		}
	}

	// 通过模板，数据渲染生成静态文件
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		ctx.Application().Logger().Error("静态文件生成失败：" + err.Error())
	}
	defer file.Close()
	template.Execute(file, product)
}

func exist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}
