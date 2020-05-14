package main

import (
	"E-commerce/common"
	"E-commerce/common/conf"
	"E-commerce/encrypt"
	"errors"
	"fmt"
	"net/http"
)

//统一验证拦截器，每个接口都需要提前验证
func Auth(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("执行验证！")

	//添加基于cookie的权限验证
	err := CheckUserInfo(r)
	if err != nil {
		return err
	}

	//return errors.New("自定义错误")
	return nil
}

func CheckUserInfo(r *http.Request) error {
	// 获取 uid , cookie
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		return errors.New("用户 UID cookie 获取失败：" + err.Error())
	}

	// 获取 UID 加密串
	signCookie, err := r.Cookie("sign")
	if err != nil {
		return errors.New("用户 sign cookie 获取失败：" + err.Error())
	}

	// 解密
	signByte, err := encrypt.DePwdCode(signCookie.Value)
	if err != nil {
		return errors.New("用户 sign cookie 解密失败，被篡改：" + err.Error())
	}
	fmt.Println("用户ID:" + uidCookie.Value + "  解密后ID：" + string(signByte))
	if uidCookie.Value == string(signByte) {
		return nil
	}
	return errors.New("身份校验失败，" + "用户ID:" + uidCookie.Value + "  解密后ID：" + string(signByte))
}

//执行正常业务逻辑
func Check(w http.ResponseWriter, r *http.Request) {
	//执行正常业务逻辑
	w.Write([]byte("检验通过，执行正常业务逻辑"))
}

func main() {

	conf.Setup()

	// 1.创建过滤器
	filter := common.NewFilter()

	// 2.注册拦截器
	filter.RegisterFilterUri("/check", Auth)

	// 3.注册路由
	http.HandleFunc("/check", filter.Handler(Check))

	// 4.启动服务
	http.ListenAndServe(":4002", nil)
}
