package services

import (
	"E-commerce/datamodels"
	"E-commerce/repositories"
)

type IProductService interface {
	GetProductByID(int64) (*datamodels.Product, error)
	GetAllProduct() ([]*datamodels.Product, error)
	DeleteProductByID(int64) bool
	InsertProduct(product *datamodels.Product) (int64, error)
	UpdateProduct(product *datamodels.Product) error
	SubNumberOne(msg *datamodels.Message) error
}

type productService struct {
	productRepository repositories.IProduct
}

func NewProductService(productRepository repositories.IProduct) IProductService {
	return &productService{productRepository: productRepository}
}

func (p *productService) GetProductByID(i int64) (*datamodels.Product, error) {
	return p.productRepository.SelectByKey(i)
}

func (p *productService) GetAllProduct() ([]*datamodels.Product, error) {
	return p.productRepository.SelectAll()
}

func (p *productService) DeleteProductByID(i int64) bool {
	return p.productRepository.Delete(i)
}

func (p *productService) InsertProduct(product *datamodels.Product) (int64, error) {
	return p.productRepository.Insert(product)
}

func (p *productService) UpdateProduct(product *datamodels.Product) error {
	return p.productRepository.Update(product)
}

func (p *productService) SubNumberOne(msg *datamodels.Message) error {
	return p.productRepository.SubProductNum(msg.ProductID, msg.UserID)
}
