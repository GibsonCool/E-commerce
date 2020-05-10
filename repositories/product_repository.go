package repositories

import (
	"E-commerce/datamodels"
	"errors"
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
	result := p.sqlDb.Create(d)
	return result.RowsAffected, result.Error
}

func (p *product) Delete(i int64) bool {
	var product datamodels.Product
	product.ID = uint(i)
	err := p.sqlDb.Delete(&product)
	return err != nil
}

func (p *product) Update(d *datamodels.Product) error {
	return p.sqlDb.Save(d).Error
}

func (p *product) SelectByKey(i int64) (*datamodels.Product, error) {
	var product datamodels.Product
	// 先查询
	e := p.sqlDb.First(&product, i).Error
	if e != nil {
		e = errors.New("商品查询不到：" + e.Error())
	}

	return &product, e
}

func (p *product) SelectAll() ([]*datamodels.Product, error) {
	var products []*datamodels.Product
	e := p.sqlDb.Find(&products).Error

	return products, e
}
