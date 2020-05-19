package main

import (
	"E-commerce/common"
	"E-commerce/common/conf"
	"E-commerce/datamodels"
	"E-commerce/encrypt"
	_0_RabbitMQ "E-commerce/rabbitmq"
	"errors"
	"fmt"
	"github.com/unknwon/com"
	"io/ioutil"
	"net/http"
	"net/url"
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
	fmt.Println("用户ID:" + uidCookie.Value + "  解密前：" + signCookie.Value + "  解密后ID：" + string(signByte))
	if uidCookie.Value == string(signByte) {
		return nil
	}
	return errors.New("身份校验失败，" + "用户ID:" + uidCookie.Value + "  解密后ID：" + string(signByte))
}

//执行正常业务逻辑
func Check(w http.ResponseWriter, r *http.Request) {
	conf.AppSetting.Logger.Info("执行check！")

	// 获取请求 URI 参数
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(query["productID"]) <= 0 {
		w.Write([]byte("false,productID 获取失败"))
		return
	}
	productStr := query["productID"][0]

	// 获取用户 cookie
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		w.Write([]byte("false,用户 uid cookie 获取失败"))
		return
	}

	// 1.分布式权限验证
	if right := accessControl.GetDistributedRight(r); right == false {
		w.Write([]byte("false,分布式权限验证出错"))
		return
	}

	// 2.获取数量权限控制，防止秒杀出现超卖
	hostUrl := "http://" + getOneIp + ":" + getOnePort + "/getOne"
	response, body, err := GetCurl(hostUrl, r)
	if err != nil {
		w.Write([]byte("false，抢单失败：" + err.Error()))
		return
	}

	// 判断数量控制接口状态
	if response.StatusCode == http.StatusOK && string(body) == "true" {
		// 创建下单消息体
		message := &datamodels.Message{
			ProductID: com.StrTo(productStr).MustInt64(),
			UserID:    com.StrTo(uidCookie.Value).MustInt64(),
		}

		// 生产消息放入 rabbitMQ
		err := rabbitMqValidate.PublishSimple(message.JsonToStr())
		if err != nil {
			w.Write([]byte("false,加入 rabbitMq 失败：" + err.Error()))
			return
		}
		w.Write([]byte("true"))
		return
	}

	w.Write([]byte("false"))
	return
}

var (
	// 设置集群地址，最好是内网IP
	hostArray      = []string{"127.0.0.1", "127.0.0.1", "127.0.0.1"}
	localHost      = "127.0.0.1"
	port           = "4002"
	consistentHash *common.ConsistentHash
	accessControl  = &AccessControl{
		sourceArray: make(map[int]interface{}),
	}
	rabbitMqValidate *_0_RabbitMQ.RabbitMQ

	// 数量控制接口服务器 IP
	getOneIp   = "127.0.0.1"
	getOnePort = "8084"
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

	// 通过一致性哈希算法，更具用户uid 获取服务器节点IP
	hostRequest, err := consistentHash.Get(uidCookie.Value)
	if err != nil {
		return false
	}

	// 判断是否为本机
	conf.AppSetting.Logger.Info("hostRequest:" + hostRequest + " localHost:" + localHost)
	if hostRequest == localHost {
		return ac.GetDataFromMap(uidCookie.Value)
	} else {
		return GetDataFromOtherMap(hostRequest, req)
	}
}

// 本机业务逻辑
func (ac *AccessControl) GetDataFromMap(value string) bool {
	//uid := com.StrTo(value).MustInt()
	//data := ac.GetNewRecord(uid)
	//
	//// 执行业务逻辑
	//if data != nil {
	//	return true
	//}
	//return false
	return true
}

// 非本机，则作为代理机，将业务请求转发到目标IP机器
func GetDataFromOtherMap(host string, req *http.Request) bool {
	hostUrl := "http://" + host + ":" + port + "/checkRight"
	response, body, err := GetCurl(hostUrl, req)
	if err != nil {
		conf.AppSetting.Logger.Error("代理请求失败：" + err.Error())
		return false
	}

	// 判断状态
	if response.StatusCode == http.StatusOK {
		return string(body) == "true"
	}
	return false
}

func CheckRight(w http.ResponseWriter, r *http.Request) {
	right := accessControl.GetDistributedRight(r)
	if !right {
		w.Write([]byte("false"))
		return
	}
	w.Write([]byte("true"))
	return
}

// 模拟转发请求
func GetCurl(hostUrl string, req *http.Request) (response *http.Response, body []byte, err error) {
	uidCookie, err := req.Cookie("uid")
	if err != nil {
		return
	}

	signCookie, err := req.Cookie("sign")
	if err != nil {
		return
	}

	// 模拟接口访问
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodGet, hostUrl, nil)
	if err != nil {
		return
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
	response, err = client.Do(request)
	if err != nil {
		return
	}

	body, err = ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	return
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

	// 获取本机IP代码
	//if ip, err := common.GetIntranceIp(); err != nil {
	//	conf.AppSetting.Logger.Error("获取本机IP失败：" + err.Error())
	//} else {
	//	localHost = ip
	//	conf.AppSetting.Logger.Info("本机IP：" + localHost)
	//}

	rabbitMqValidate = _0_RabbitMQ.NewRabbitMQSimple(_0_RabbitMQ.TestQueueName)
	defer rabbitMqValidate.Destory()

	// 1.创建过滤器
	filter := common.NewFilter()

	// 2.注册拦截器
	filter.RegisterFilterUri("/check", Auth)
	filter.RegisterFilterUri("/checkRight", Auth)

	// 3.注册路由
	http.HandleFunc("/check", filter.Handler(Check))
	http.HandleFunc("/checkRight", filter.Handler(CheckRight))

	// 4.启动服务
	http.ListenAndServe(":"+port, nil)
}
