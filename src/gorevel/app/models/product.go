package models

import (
	"time"
	"unicode/utf8"

	"github.com/robfig/revel"
)

type Product struct {
	Id          int64
	Name        string
	Site        string
	Author      string
	Repository  string
	Description string `xorm:"text"`
	Image       string
	User        User      `xorm:"user_id bigint"`
	Created     time.Time `xorm:"created"`
	Updated     time.Time `xorm:"updated"`
}

func (product Product) Validate(v *revel.Validation) {
	v.Required(product.Name).Message("不能为空")
	v.Required(product.Author).Message("不能为空")
	v.Required(product.Description).Message("不能为空")

	if utf8.RuneCountInString(product.Name) > 20 {
		err := &revel.ValidationError{
			Message: "最多20个字",
			Key:     "product.Name",
		}
		v.Errors = append(v.Errors, err)
	}
}

func (p Product) GetImage() string {
	return QiniuDomain + "/" + p.Image
}
