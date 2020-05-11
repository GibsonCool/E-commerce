package repositories

import (
	"E-commerce/datamodels"
	"errors"
	"github.com/jinzhu/gorm"
)

type IOrder interface {
	Insert(*datamodels.Order) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Order) error
	SelectByKey(int64) (*datamodels.Order, error)
	SelectAll() ([]*datamodels.Order, error)
	SelectAllWithInfo() ([]datamodels.OrderInfo, error)
}

type order struct {
	sqlDb *gorm.DB
}

func (o *order) Insert(d *datamodels.Order) (int64, error) {
	e := o.sqlDb.Create(d).Error
	if e != nil {
		return 0, errors.New("插入订单失败：" + e.Error())
	}
	return int64(d.ID), e
}

func (o *order) Delete(i int64) bool {
	panic("implement me")
}

func (o *order) Update(d *datamodels.Order) error {
	panic("implement me")
}

func (o *order) SelectByKey(i int64) (*datamodels.Order, error) {
	panic("implement me")
}

func (o *order) SelectAll() ([]*datamodels.Order, error) {
	panic("implement me")
}

func (o *order) SelectAllWithInfo() (results []datamodels.OrderInfo, err error) {
	err = o.sqlDb.Table("order").
		Select("order.id,user.user_name,product.product_name,order.order_status").
		Joins("left join product on order.product_id=product.id").
		Joins("left join user on order.user_id=user.id").
		Scan(&results).Error
	if err != nil {
		err = errors.New("查询所有订单信息失败：" + err.Error())
	}
	return results, err
}

func NewOrder(sqlDb *gorm.DB) *order {
	return &order{sqlDb: sqlDb}
}
