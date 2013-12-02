package controllers

import (
	"fmt"
	"github.com/coocood/qbs"
	"github.com/robfig/revel"
	"gorevel/app/models"
	"gorevel/app/routes"
)

type Admin struct {
	Application
}

func (c *Admin) Index() revel.Result {
	return c.Render()
}

func (c *Admin) ListUser() revel.Result {
	var users []*models.User
	c.q.FindAll(&users)

	return c.Render(users)
}

func (c *Admin) DeleteUser(id int64) revel.Result {
	user := new(models.User)
	user.Id = id
	c.q.Delete(user)

	return c.RenderJson([]byte("true"))
}

func (c *Admin) ListCategory() revel.Result {
	categories := getCategories(c.q)

	return c.Render(categories)
}

func (c *Admin) DeleteCategory(id int64) revel.Result {
	category := new(models.Category)
	category.Id = id
	c.q.Delete(category)

	return c.RenderJson([]byte("true"))
}

func (c *Admin) NewCategory() revel.Result {
	title := "新建分类"
	return c.Render(title)
}

func (c *Admin) NewCategoryPost(category models.Category) revel.Result {
	category.Validate(c.q, c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Admin.NewCategory())
	}

	if !category.Save(c.q) {
		c.Flash.Error("添加分类失败")
	}

	return c.Redirect(routes.Admin.ListCategory())
}

func (c *Admin) EditCategory(id int64) revel.Result {
	title := "编辑分类"

	category := findCategoryById(c.q, id)
	if category.Id == 0 {
		return c.NotFound("分类不存在")
	}

	c.Render(title, category)

	return c.RenderTemplate("admin/NewCategory.html")
}

func (c *Admin) EditCategoryPost(id int64, category models.Category) revel.Result {
	category.Id = id
	category.Validate(c.q, c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Admin.NewCategory())
	}

	if !category.Save(c.q) {
		c.Flash.Error("编辑分类失败")
	}

	return c.Redirect(routes.Admin.ListCategory())
}

func getCategories(q *qbs.Qbs) []*models.Category {
	var categories []*models.Category
	if err := q.FindAll(&categories); err != nil {
		fmt.Println(err)
	}

	return categories
}

func findCategoryById(q *qbs.Qbs, id int64) *models.Category {
	category := new(models.Category)
	q.WhereEqual("id", id).Find(category)

	return category
}
