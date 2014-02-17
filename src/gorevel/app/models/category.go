package models

import (
	"time"

	"github.com/robfig/revel"
)

type Category struct {
	Id      int64
	Name    string
	Intro   string
	Created time.Time `xorm:"created"`
}

func (category Category) Validate(v *revel.Validation) {
	v.Required(category.Name).Message("请输入名称")

	if category.HasName() {
		err := &revel.ValidationError{
			Message: "名称已存在",
			Key:     "category.Name",
		}
		v.Errors = append(v.Errors, err)
	}
}

func (c Category) HasName() bool {
	var category Category
	if c.Id > 0 {
		Engine.Where("name = ? AND id != ?", c.Name, c.Id).Get(&category)
	} else {
		Engine.Where("name = ?", c.Name).Get(&category)
	}

	return category.Id > 0
}
