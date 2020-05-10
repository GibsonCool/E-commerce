package main

import (
	"E-commerce/backend/web/controllers"
	"E-commerce/common/conf"
	"E-commerce/common/datasourse"
	"E-commerce/repositories"
	"E-commerce/services"
	"context"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

func main() {
	conf.Setup()

	app := newApp()
	// 路由设置
	mvcHandler(app)

	app.Run(
		iris.Addr(":4000"),
		iris.WithoutServerError(iris.ErrServerClosed), //无服务错误提示
		iris.WithOptimizations,                        //让程序自身尽可能的优化
	)
}

func newApp() *iris.Application {
	// 1、创建 iris 实例
	app := iris.New()

	// 2.设置日志级别，开发阶段为 debug
	conf.AppSetting.Logger = app.Logger().SetLevel("debug")

	// 3.注册静态资源
	app.StaticWeb("/assets", "./backend/web/assets")

	// 4.注册模板
	template := iris.HTML("./backend/web/views", ".html").
		Layout("shared/layout.html").
		Reload(true)
	app.RegisterView(template)

	// 5.设置异常出错处理
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("Message", ctx.Values().GetStringDefault("Message", "访问的页面出错！"+ctx.Path()))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})

	return app
}

func mvcHandler(app *iris.Application) {

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// 商品管理控制器
	productRepository := repositories.NewProduct(datasourse.GetInstance())
	productService := services.NewProductService(productRepository)
	productGroup := mvc.New(app.Party("/product"))
	productGroup.Register(
		ctx,
		productService,
	)
	productGroup.Handle(new(controllers.ProductController))

	// 订单管理控制器
	orderRepository := repositories.NewOrder(datasourse.GetInstance())
	orderService := services.NewOrderService(orderRepository)
	orderGroup := mvc.New(app.Party("/order"))
	orderGroup.Register(ctx, orderService)
	orderGroup.Handle(new(controllers.OrderController))

}
