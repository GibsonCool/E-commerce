package main

import (
	"E-commerce/common"
	"E-commerce/common/conf"
	"E-commerce/encrypt"
	"errors"
	"fmt"
	"github.com/unknwon/com"
	"io/ioutil"
	"net/http"
	"sync"
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

var (
	hostArray      = []string{"127.0.0.1", "127.0.0.1", "127.0.0.1"}
	localHost      = "127.0.0.1"
	port           = "8081"
	consistentHash *common.ConsistentHash
	accessControl  = &AccessControl{
		sourceArray: make(map[int]interface{}),
	}
)

// 存放访问控制数据信息
type AccessControl struct {
	// 用于存放用户想要存放的信息
	sourceArray map[int]interface{}
	sync.RWMutex
}

// 获取数据
func (ac *AccessControl) GetNewRecord(uid int) interface{} {
	ac.RLock()
	defer ac.RUnlock()
	return ac.sourceArray[uid]
}

// 设置数据
func (ac *AccessControl) SetNewRecord(uid int) {
	ac.Lock()
	defer ac.Unlock()
	ac.sourceArray[uid] = "test hello world"
}

func (ac *AccessControl) GetDistributedRight(req *http.Request) bool {
	// 获取用户 UID
	uidCookie, err := req.Cookie("uid")
	if err != nil {
		return false
	}

	// 通过一致性哈希算法，更具用户ID 获取服务器节点IP
	hostRequest, err := consistentHash.Get(uidCookie.Value)
	if err != nil {
		return false
	}

	// 判断是否为本机
	if hostRequest == localHost {
		return ac.GetDataFromMap(uidCookie.Value)
	} else {
		return ac.GetDataFromOtherMap(hostRequest, req)
	}
}

func (ac *AccessControl) GetDataFromMap(value string) bool {
	uid := com.StrTo(value).MustInt()
	data := ac.GetNewRecord(uid)

	// 执行业务逻辑
	if data != nil {
		return true
	}
	return false
}

func (ac *AccessControl) GetDataFromOtherMap(host string, req *http.Request) bool {
	uidCookie, err := req.Cookie("uid")
	if err != nil {
		return false
	}

	signCookie, err := req.Cookie("sign")
	if err != nil {
		return false
	}

	// 模拟接口访问
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodGet, "http://"+host+":"+port+"check", nil)
	if err != nil {
		return false
	}

	// 手动指定，排除多余 cookies
	cookieUid := &http.Cookie{
		Name:  "uid",
		Value: uidCookie.Value,
		Path:  "/",
	}

	cookieSign := &http.Cookie{
		Name:  "sign",
		Value: signCookie.Value,
		Path:  "/",
	}
	// 添加 cookie 到模拟请求中
	request.AddCookie(cookieSign)
	request.AddCookie(cookieUid)
	// 发起请求，获取返回结果
	response, err := client.Do(request)
	if err != nil {
		return false
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false
	}

	// 判断状态
	if response.StatusCode == http.StatusOK {
		return string(body) == "true"
	}
	return false
}

func main() {

	conf.Setup()

	// 负载均衡器设置
	// 采用一致性hash算法
	consistentHash = common.NewConsistentHash()
	// 添加节点
	for _, value := range hostArray {
		consistentHash.Add(value)
	}

	// 1.创建过滤器
	filter := common.NewFilter()

	// 2.注册拦截器
	filter.RegisterFilterUri("/check", Auth)

	// 3.注册路由
	http.HandleFunc("/check", filter.Handler(Check))

	// 4.启动服务
	http.ListenAndServe(":4002", nil)
}
