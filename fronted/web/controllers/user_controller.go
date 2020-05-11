package controllers

import (
	"E-commerce/datamodels"
	"E-commerce/services"
	"E-commerce/tool"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"strconv"
)

type UserController struct {
	Ctx         iris.Context
	UserService services.IUserService
	Session     *sessions.Sessions
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

	uc.Session.Start(uc.Ctx).Set("userID", strconv.FormatInt(int64(user.ID), 10))

	return mvc.Response{
		Path: "/product/",
	}
}
