package controllers

import (
	"github.com/robfig/revel"
)

func init() {
	revel.OnAppStart(Init)

	revel.InterceptMethod((*Qbs).Begin, revel.BEFORE)
	revel.InterceptMethod((*Application).inject, revel.BEFORE)
	revel.InterceptMethod((*Qbs).End, revel.AFTER)
}
