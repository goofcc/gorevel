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
	products := getProducts()

	return c.Render(products)
}

func (c Product) New() revel.Result {
	title := "提交案例"

	return c.Render(title)
}

func (c Product) NewPost(product models.Product) revel.Result {
	product.Validate(c.Validation)

	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()
		if ok := checkImageExt(c.Validation, &file, header, "image"); ok {
			fileName := uuidFileName(header.Filename)
			err, ret := qiniuUploadImage(&file, fileName)
			if err != nil {
				c.Flash.Error("上传到七牛出错，请检查七牛配置。")
				return c.Redirect(routes.Product.New())
			} else {
				product.Image = ret.Key
			}
		}
	} else {
		err := &revel.ValidationError{
			Message: "字段不能为空",
			Key:     "image",
		}
		c.Validation.Errors = append(c.Validation.Errors, err)
	}

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Product.New())
	}

	product.User = models.User{Id: c.user().Id}

	aff, _ := engine.Insert(&product)
	if aff > 0 {
		c.Flash.Success("提交案例成功")
		cache.Delete("products")
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
	var tmp models.Product
	has, _ := engine.Id(id).Get(&tmp)
	if !has {
		return c.NotFound("案例不存在")
	}

	product.Validate(c.Validation)

	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()
		if ok := checkImageExt(c.Validation, &file, header, "image"); ok {
			fileName := uuidFileName(header.Filename)
			err, ret := qiniuUploadImage(&file, fileName)
			if err != nil {
				c.Flash.Error("上传到七牛出错，请检查七牛配置。")
				return c.Redirect(routes.Product.Edit(id))
			} else {
				if tmp.Image != "" {
					qiniuDeleteImage(tmp.Image)
				}
				product.Image = ret.Key
			}
		}
	} else {
		product.Image = tmp.Image
	}

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Product.Edit(id))
	}

	// 强制更新允许空值的字段
	aff, _ := engine.Id(id).MustCols("site", "repository").Update(&product)
	if aff > 0 {
		c.Flash.Success("编辑案例成功")
		cache.Delete("products")
	} else {
		c.Flash.Error("编辑案例失败")
		return c.Redirect(routes.Product.Edit(id))
	}

	return c.Redirect(routes.Product.Index())
}

func getProducts() []models.Product {
	var products []models.Product
	if err := cache.Get("products", &products); err != nil {
		engine.Find(&products)
		go cache.Set("products", products, cache.FOREVER)
	}

	return products
}
