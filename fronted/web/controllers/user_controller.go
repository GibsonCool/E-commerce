package controllers

import (
	"E-commerce/datamodels"
	"E-commerce/encrypt"
	"E-commerce/services"
	"E-commerce/tool"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"github.com/unknwon/com"
	"strconv"
)

type UserController struct {
	Ctx         iris.Context
	UserService services.IUserService
	Session     *sessions.Session
}

// /user/register
// get
func (uc *UserController) GetRegister() mvc.View {
	return mvc.View{
		Name: "user/register.html",
	}
}

// /user/register
// post
func (uc *UserController) PostRegister() {
	user := &datamodels.User{
		NickName: uc.Ctx.FormValue("nickName"),
		UserName: uc.Ctx.FormValue("userName"),
		HashPwd:  uc.Ctx.FormValue("password"),
	}

	_, err := uc.UserService.AddUser(user)
	if err != nil {
		uc.Ctx.Redirect("/user/error")
		return
	}
	uc.Ctx.Redirect("/user/login")
	return
}

// /user/login
// get
func (uc *UserController) GetLogin() mvc.View {
	return mvc.View{
		Name: "user/login.html",
	}
}

// /user/login
// post
func (uc *UserController) PostLogin() mvc.Result {
	userName := uc.Ctx.FormValue("userName")
	pwd := uc.Ctx.FormValue("password")

	// 查询用户校验是否正确
	user, isOk := uc.UserService.IsPwdSuccess(userName, pwd)
	if !isOk {
		return mvc.Response{
			Path: "/user/login",
		}
	}

	// 将用户ID写入 cookie 中
	tool.GlobalCookie(uc.Ctx, "uid", strconv.FormatInt(int64(user.ID), 10))

	// 取消 session减小服务器压力，依然使用 cookie,进行加密处理，
	//uc.Session.Set("userID", strconv.FormatInt(int64(user.ID), 10))
	uidByte := []byte(com.ToStr(user.ID))
	uidStr, err := encrypt.EnPwdCode(uidByte)
	uc.Ctx.Application().Logger().Error("uidStr:" + uidStr)
	if err != nil {
		uc.Ctx.Application().Logger().Error(err.Error())
	}
	// 写入用户浏览及
	tool.GlobalCookie(uc.Ctx, "sign", uidStr)

	return mvc.Response{
		Path: "/product/",
	}
}
