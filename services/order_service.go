package services

import (
	"E-commerce/datamodels"
	"E-commerce/repositories"
)

type IOrderService interface {
	GetOrderByID(int64) (*datamodels.Order, error)
	DeleteOrderByID(int64) bool
	UpdateOrder(*datamodels.Order) error
	InsertOrder(*datamodels.Order) (int64, error)
	GetAllOrder() ([]*datamodels.Order, error)
	GetAllOrderInfo() ([]datamodels.OrderInfo, error)
}

type orderService struct {
	OrderRepository repositories.IOrder
}

func NewOrderService(orderRepository repositories.IOrder) *orderService {
	return &orderService{OrderRepository: orderRepository}
}

func (o *orderService) GetOrderByID(i int64) (*datamodels.Order, error) {
	return o.OrderRepository.SelectByKey(i)
}

func (o *orderService) DeleteOrderByID(i int64) bool {
	return o.OrderRepository.Delete(i)
}

func (o *orderService) UpdateOrder(order *datamodels.Order) error {
	return o.OrderRepository.Update(order)
}

func (o *orderService) InsertOrder(order *datamodels.Order) (int64, error) {
	return o.OrderRepository.Insert(order)
}

func (o *orderService) GetAllOrder() ([]*datamodels.Order, error) {
	return o.OrderRepository.SelectAll()
}

func (o *orderService) GetAllOrderInfo() ([]datamodels.OrderInfo, error) {
	return o.OrderRepository.SelectAllWithInfo()
}
