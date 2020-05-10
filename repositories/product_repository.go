package repositories

import (
	"E-commerce/datamodels"
	"github.com/jinzhu/gorm"
)

type IProduct interface {
	Insert(*datamodels.Product) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Product) error
	SelectByKey(int64) (*datamodels.Product, error)
	SelectAll() ([]*datamodels.Product, error)
}

type product struct {
	sqlDb *gorm.DB
}

func NewProduct(sqlDb *gorm.DB) *product {
	return &product{sqlDb: sqlDb}
}

func (p *product) Insert(d *datamodels.Product) (int64, error) {
	panic("implement me")
}

func (p *product) Delete(i int64) bool {
	panic("implement me")
}

func (p *product) Update(d *datamodels.Product) error {
	panic("implement me")
}

func (p *product) SelectByKey(i int64) (*datamodels.Product, error) {
	panic("implement me")
}

func (p *product) SelectAll() ([]*datamodels.Product, error) {
	var products []*datamodels.Product
	e := p.sqlDb.Find(&products).Error

	return products,e
}
