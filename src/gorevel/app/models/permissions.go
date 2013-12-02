package models

import (
	"fmt"
	"github.com/coocood/qbs"
)

type Permissions struct {
	Id     int64
	UserId int64 `qbs:"notnull"`
	Perm   int   `qbs:"notnull"`
}

func (p *Permissions) Save(q *qbs.Qbs) bool {
	_, err := q.Save(p)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
