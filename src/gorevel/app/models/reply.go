package models

import (
	"time"
)

type Reply struct {
	Id      int64
	Topic   Topic     `xorm:"topic_id bigint"`
	User    User      `xorm:"user_id bigint"`
	Content string    `xorm:"text"`
	Created time.Time `xorm:"created"`
}
