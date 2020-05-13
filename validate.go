package main

import (
	"E-commerce/common"
	"errors"
	"fmt"
	"net/http"
)

//统一验证拦截器，每个接口都需要提前验证
func Auth(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("执行验证！")

	//添加基于cookie的权限验证
	//err := CheckUserInfo(r)
	//if err != nil {
	//	return err
	//}

	return errors.New("自定义错误")
	//return nil
}

//执行正常业务逻辑
func Check(w http.ResponseWriter, r *http.Request) {
	//执行正常业务逻辑
	fmt.Println("执行check！")
}

func main() {

	// 1.创建过滤器
	filter := common.NewFilter()

	// 2.注册拦截器
	filter.RegisterFilterUri("/check", Auth)

	// 3.注册路由
	http.HandleFunc("/check", filter.Handler(Check))

	// 4.启动服务
	http.ListenAndServe(":4002", nil)
}
