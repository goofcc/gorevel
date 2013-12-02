package models

import (
	"fmt"
	"github.com/coocood/qbs"
	"github.com/robfig/revel"
	"strings"
	"time"
)

type Topic struct {
	Id         int64
	Title      string `qbs:"size:255,notnull"`
	Content    string `qbs:"notnull"`
	CategoryId int64  `qbs:"notnull"`
	Category   *Category
	UserId     int64 `qbs:"notnull"`
	User       *User
	Hits       int
	Replies    int
	Good       bool
	Created    time.Time
	Updated    time.Time
}

func (topic *Topic) Validate(v *revel.Validation) {
	v.Required(topic.Title).Message("请输入标题")
	v.MaxSize(topic.Title, 105).Message("最多35个字")
	v.Required(topic.Category).Message("请选择分类")
	v.Required(topic.Content).Message("帖子内容不能为空")
}

func (t *Topic) Save(q *qbs.Qbs) bool {
	_, err := q.Save(t)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func GetTopics(q *qbs.Qbs, page int, column string, value interface{}, order string, url string) ([]*Topic, *Pagination) {
	page -= 1
	if page < 0 {
		page = 0
	}

	var topics []*Topic
	var rows int64
	if column == "" {
		rows = q.Count("topic")
		err := q.OmitFields("Content").OrderByDesc(order).
			Limit(ItemsPerPage).Offset(page * ItemsPerPage).FindAll(&topics)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		rows = q.WhereEqual(column, value).Count("topic")
		err := q.WhereEqual(column, value).
			OmitFields("Content").OrderByDesc(order).
			Limit(ItemsPerPage).Offset(page * ItemsPerPage).FindAll(&topics)
		if err != nil {
			fmt.Println(err)
		}
	}

	url = url[:strings.Index(url, "=")+1]
	pagination := NewPagination(page, int(rows), url)

	return topics, pagination
}
