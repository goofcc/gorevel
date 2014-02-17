package models

import (
	"time"

	"github.com/robfig/revel"
)

type Topic struct {
	Id       int64
	Title    string
	Content  string   `xorm:"text"`
	Category Category `xorm:"category_id"`
	User     User     `xorm:"user_id"`
	Hits     int
	Replies  int
	Good     bool
	Created  time.Time `xorm:"created"`
	Updated  time.Time `xorm:"updated"`
}

func (topic Topic) Validate(v *revel.Validation) {
	v.Required(topic.Title).Message("请输入标题")
	v.MaxSize(topic.Title, 105).Message("最多35个字")
	v.Required(topic.Category).Message("请选择分类")
	v.Required(topic.Content).Message("帖子内容不能为空")
}
