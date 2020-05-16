package conf

import (
	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"gopkg.in/ini.v1"
	"log"
	"time"
)

type App struct {
	JwtSecret string
	PageSize  int
	AesKey    string
	Logger    *golog.Logger
}

type DataBase struct {
	Type        string
	SqlDSN      string
	MaxIdle     int
	MaxOpen     int
	IdleTimeout time.Duration
}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
	NetWork     string
	Prefix      string
}

var (
	AppSetting      = &App{}
	DataBaseSetting = &DataBase{}
	RedisSetting    = &Redis{}
)

/*
	编写(App、Server、DataBase)与 app.ini 一直的结构体
	使用 MapTo 将配置项映射到结构上面定义的结构体上
	对一些特殊设置的配置项进行在赋值
*/
var cfg *ini.File

func Setup() {
	var e error
	cfg, e = ini.Load("./common/conf/app.ini")
	if e != nil {
		log.Fatalf("Fail to parse 'conf/app.in' : %v", e)
	}

	mapTo("app", AppSetting)
	mapTo("database", DataBaseSetting)
	mapTo("redis", RedisSetting)
	AppSetting.Logger = iris.New().Logger()

	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
}

func mapTo(section string, v interface{}) {
	e := cfg.Section(section).MapTo(v)
	if e != nil {
		log.Fatalf("cfg.MapTo %sSetting err: %v", section, e)
	}
}
