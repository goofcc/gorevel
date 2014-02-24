package models

import (
	"time"
	"unicode/utf8"

	"github.com/robfig/revel"
)

type Topic struct {
	Id       int64
	Title    string
	Content  string   `xorm:"text"`
	Category Category `xorm:"category_id bigint"`
	User     User     `xorm:"user_id bigint"`
	Hits     int
	Replies  int
	Good     bool
	Created  time.Time `xorm:"created"`
	Updated  time.Time `xorm:"updated"`
}

func (topic Topic) Validate(v *revel.Validation) {
	v.Required(topic.Title).Message("请输入标题")
	if utf8.RuneCountInString(topic.Title) > 35 {
		err := &revel.ValidationError{
			Message: "最多35个字",
			Key:     "topic.Title",
		}
		v.Errors = append(v.Errors, err)
	}
	v.Required(topic.Category).Message("请选择分类")
	v.Required(topic.Content).Message("帖子内容不能为空")
}
