package models

import (
	"fmt"
	"github.com/coocood/qbs"
	"github.com/robfig/revel"
	"time"
)

type Category struct {
	Id      int64
	Name    string `qbs:"size:32,unique,notnull"`
	Intro   string `qbs:"size:255"`
	Created time.Time
}

func (category *Category) Validate(q *qbs.Qbs, v *revel.Validation) {
	valid := v.Required(category.Name).Message("请输入名称")

	if valid.Ok {
		if category.HasName(q) {
			err := &revel.ValidationError{
				Message: "该名称已存在",
				Key:     "category.Name",
			}
			valid.Error = err
			valid.Ok = false

			v.Errors = append(v.Errors, err)
		}
	}
}

func (c *Category) HasName(q *qbs.Qbs) bool {
	category := new(Category)
	condition := qbs.NewCondition("name = ?", c.Name)
	if c.Id > 0 {
		condition = qbs.NewCondition("name = ?", c.Name).And("id != ?", c.Id)
	}
	q.Condition(condition).Find(category)

	return category.Id > 0
}

func (c *Category) Save(q *qbs.Qbs) bool {
	_, err := q.Save(c)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
