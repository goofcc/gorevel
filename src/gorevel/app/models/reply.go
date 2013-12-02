package models

import (
	"fmt"
	"github.com/coocood/qbs"
	"time"
)

type Reply struct {
	Id      int64
	TopicId int64 `qbs:"notnull"`
	Topic   *Topic
	UserId  int64 `qbs:"notnull"`
	User    *User
	Content string `qbs:"notnull"`
	Created time.Time
}

func (r *Reply) Save(q *qbs.Qbs) bool {
	_, err := q.Save(r)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
