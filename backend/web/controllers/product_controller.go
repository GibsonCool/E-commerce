package controllers

import (
	"E-commerce/common/conf"
	"E-commerce/datamodels"
	"E-commerce/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/unknwon/com"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
}

// 获取商品列表
// url:  /product/all
// type: get
func (p *ProductController) GetAll() mvc.View {
	product, err := p.ProductService.GetAllProduct()
	if err != nil {
		conf.AppSetting.Logger.Error(err.Error())
	}

	return mvc.View{
		Name: "product/view.html",
		Data: iris.Map{
			"productArray": product,
		},
	}
}

//修改商品
// url:  /product/update
// type: post
func (p *ProductController) PostUpdate() {
	product := &datamodels.Product{
		ProductName:  p.Ctx.FormValue("ProductName"),
		ProductNum:   com.StrTo(p.Ctx.FormValue("ProductNum")).MustInt64(),
		ProductImage: p.Ctx.FormValue("ProductImage"),
		ProductUrl:   p.Ctx.FormValue("ProductUrl"),
	}
	product.ID = uint(com.StrTo(p.Ctx.FormValue("ID")).MustInt())
	p.Ctx.Application().Logger().Error(product)

	err := p.ProductService.UpdateProduct(product)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	p.Ctx.Redirect("/product/all")
}

// 添加商品页面
// url:  /product/add
// type: get
func (p *ProductController) GetAdd() mvc.View {
	return mvc.View{
		Name: "product/add.html",
	}
}

// 添加商品
// url:  /product/all
// type: post
func (p *ProductController) PostAdd() {
	product := &datamodels.Product{
		ProductName:  p.Ctx.FormValue("ProductName"),
		ProductNum:   com.StrTo(p.Ctx.FormValue("ProductNum")).MustInt64(),
		ProductImage: p.Ctx.FormValue("ProductImage"),
		ProductUrl:   p.Ctx.FormValue("ProductUrl"),
	}
	rows, err := p.ProductService.InsertProduct(product)
	p.Ctx.Application().Logger().Error("插入数据：" + com.ToStr(rows))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	p.Ctx.Redirect("/product/all")
}

// 商品修改页面
// url:  /product/manager
// type: get
func (p *ProductController) GetManager() mvc.View {
	idString := p.Ctx.URLParam("id")
	id := com.StrTo(idString).MustInt64()

	product, err := p.ProductService.GetProductByID(id)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	return mvc.View{
		Name: "product/manager.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

// 删除商品
// url:  /product/delete
// type: get
func (p *ProductController) GetDelete() {
	idString := p.Ctx.URLParam("id")
	id := com.StrTo(idString).MustInt64()
	isOk := p.ProductService.DeleteProductByID(id)
	if isOk {
		p.Ctx.Application().Logger().Debug("删除商品成功，ID为：" + idString)
	} else {
		p.Ctx.Application().Logger().Debug("删除商品失败，ID为：" + idString)
	}
	p.Ctx.Redirect("/product/all")
}
