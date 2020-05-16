package datamodels

import (
	"E-commerce/common/conf"
	"encoding/json"
)

// 商品下单用于消息队列传输的消息
type Message struct {
	ProductID int64
	UserID    int64
}

func NewMessage(productID int64, userID int64) *Message {
	return &Message{ProductID: productID, UserID: userID}
}

func (m *Message) JsonToStr() string {
	bytes, err := json.Marshal(m)
	if err != nil {
		conf.AppSetting.Logger.Error("json 解析出错：" + err.Error())
	}

	return string(bytes)
}

func (m *Message) StrToJson(dataStr []byte) *Message {
	if err := json.Unmarshal(dataStr, m); err != nil {
		conf.AppSetting.Logger.Error("json 转换出错：" + err.Error())
	}
	return m
}
