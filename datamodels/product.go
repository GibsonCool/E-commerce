package datamodels

import (
	"github.com/jinzhu/gorm"
	"github.com/unknwon/com"
)

type Product struct {
	gorm.Model
	ProductName  string
	ProductNum   int64
	ProductImage string
	ProductUrl   string
}

func (p *Product) ToString() string {
	return "productNum:" + com.ToStr(p.ProductNum)
}
