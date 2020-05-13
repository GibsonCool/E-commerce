package middleware

import (
	"E-commerce/encrypt"
	"github.com/kataras/iris"
)

func AuthConProduct(ctx iris.Context) {
	uid := ctx.GetCookie("sign")
	if uid == "" {
		ctx.Application().Logger().Debug("必须先登录!")
		ctx.Redirect("/user/login")
		return
	}
	code, err := encrypt.DePwdCode(uid)
	if err != nil {
		ctx.Application().Logger().Error("登录校验失败：" + err.Error())
		return
	}
	ctx.Application().Logger().Debug("sign:" + string(code))
	ctx.Next()
}
