package main

import (
	"E-commerce/common/conf"
	"E-commerce/common/datasourse"
	"E-commerce/fronted/web/controllers"
	"E-commerce/repositories"
	"E-commerce/services"
	"context"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"time"
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

	sess := sessions.New(sessions.Config{
		Cookie:  "AdminCookie",
		Expires: 600 * time.Minute,
	})
	// 设置 session 使用 redis 来保存信息
	sess.UseDatabase(datasourse.GetRedisInstance())

	// 用户管理控制器
	userRepository := repositories.NewUserRepository(datasourse.GetMysqlInstance())
	userService := services.NewUserService(userRepository)
	userGroup := mvc.New(app.Party("/user"))
	userGroup.Register(ctx, userService, sess)
	userGroup.Handle(new(controllers.UserController))

	// 商品管理控制器
	productRepository := repositories.NewProduct(datasourse.GetMysqlInstance())
	productService := services.NewProductService(productRepository)
	productGroup := mvc.New(app.Party("/product"))
	productGroup.Register(ctx, &productService)
	productGroup.Handle(new(controllers.ProductController))
}
