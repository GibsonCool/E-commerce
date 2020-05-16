package main

import (
	"E-commerce/common/conf"
	"E-commerce/common/datasourse"
	"E-commerce/fronted/middleware"
	"E-commerce/fronted/web/controllers"
	_0_RabbitMQ "E-commerce/rabbitmq"
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
		iris.Addr(":4001"),
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
	app.StaticWeb("/public", "./fronted/web/public")
	// 3.1 注册静态模板生成页面
	app.StaticWeb("/html", "./fronted/web/htmlProductShow")

	// 4.注册模板
	template := iris.HTML("./fronted/web/views", ".html").
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

	//sess := sessions.New(sessions.Config{
	//	Cookie:  "AdminCookie",
	//	Expires: 600 * time.Minute,
	//})
	//// 设置 session 使用 redis 来保存信息
	//sess.UseDatabase(datasourse.GetRedisInstance())

	// 用户管理控制器
	userRepository := repositories.NewUserRepository(datasourse.GetMysqlInstance())
	userService := services.NewUserService(userRepository)
	userGroup := mvc.New(app.Party("/user"))

	// 使用 cookie 代替 session 后可以取消使用
	//userGroup.Register(ctx, userService, sess.Start)
	userGroup.Register(ctx, userService)

	userGroup.Handle(new(controllers.UserController))

	// 商品管理控制器
	productRepository := repositories.NewProduct(datasourse.GetMysqlInstance())
	productService := services.NewProductService(productRepository)

	orderRepository := repositories.NewOrder(datasourse.GetMysqlInstance())
	orderService := services.NewOrderService(orderRepository)

	productParty := app.Party("/product")
	// 使用中间件，进行登录校验
	productParty.Use(middleware.AuthConProduct)

	rabbitMQSimple := _0_RabbitMQ.NewRabbitMQSimple(_0_RabbitMQ.TestQueueName)

	productGroup := mvc.New(productParty)
	productGroup.Register(ctx, productService, orderService, rabbitMQSimple)
	productGroup.Handle(new(controllers.ProductController))
}
