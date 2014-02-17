package controllers

import (
	"strconv"

	"github.com/robfig/revel"

	"gorevel/app/models"
)

func init() {
	revel.OnAppStart(Init)
	revel.InterceptMethod((*Application).injector, revel.BEFORE)
}

func Init() {
	engine = models.Engine

	revel.TemplateFuncs["eqis"] = func(a int64, b string) bool {
		s := strconv.FormatInt(a, 10)
		return s == b
	}

	UPLOAD_PATH = revel.BasePath + "/public/upload/"
}
