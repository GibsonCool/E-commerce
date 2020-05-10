package datamodels

import "github.com/jinzhu/gorm"

type Order struct {
	gorm.Model
	UserId      int64
	ProductId   int64
	OrderStatus int64
}

type OrderInfo struct {
	ID          int64
	UserName    string
	ProductName string
	OrderStatus string
}

// 订单状态
const (
	OrderWait = iota
	OrderSuccess
	OrderFailed
)