package models

import (
	"time"
)

type Reply struct {
	Id      int64
	Topic   Topic     `xorm:"topic_id"`
	User    User      `xorm:"user_id"`
	Content string    `xorm:"text"`
	Created time.Time `xorm:"created"`
}
