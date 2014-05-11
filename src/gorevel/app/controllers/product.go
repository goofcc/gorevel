package controllers

import (
	"github.com/robfig/revel"
	"github.com/robfig/revel/cache"

	"gorevel/app/models"
	"gorevel/app/routes"
)

type Product struct {
	Application
}

func (c Product) Index() revel.Result {
	var products []models.Product
	engine.Find(&products)

	return c.Render(products)
}

func (c Product) New() revel.Result {
	title := "提交案例"

	return c.Render(title)
}

func (c Product) NewPost(product models.Product) revel.Result {
	product.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Product.New())
	}

	product.User = models.User{Id: c.user().Id}

	aff, _ := engine.Insert(&product)
	if aff > 0 {
		c.Flash.Success("提交案例成功")
		cache.Flush()
	} else {
		c.Flash.Error("提交案例失败")
	}

	return c.Redirect(routes.Product.Index())
}

func (c Product) Edit(id int64) revel.Result {
	var product models.Product
	has, _ := engine.Id(id).Get(&product)
	if !has {
		return c.NotFound("案例不存在")
	}

	c.bindVars(Vars{
		"title":   "编辑案例",
		"product": product,
	})

	return c.RenderTemplate("product/New.html")
}

func (c Product) EditPost(id int64, product models.Product) revel.Result {
	product.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Product.Edit(id))
	}

	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()
		if ok := checkFileExt(c.Validation, header, imageExts, "image", "Only image"); ok {
			fileName := uuidFileName(header.Filename)
			err, ret := qiniuUploadImage(&file, fileName)
			if err != nil {
				c.Flash.Error("上传到七牛出错，请检查七牛配置。")
				return c.Redirect(routes.User.Edit())
			} else {
				if product.Image != "" {
					qiniuDeleteImage(product.Image)
				}
				product.Image = ret.Key
			}
		}
	}

	aff, _ := engine.Id(id).Cols("name", "site", "author", "repository", "description", "image").Update(&product)
	if aff > 0 {
		c.Flash.Success("编辑案例成功")
		cache.Flush()
	} else {
		c.Flash.Error("编辑案例失败")
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Product.Edit(id))
	}

	return c.Redirect(routes.Product.Index())
}
