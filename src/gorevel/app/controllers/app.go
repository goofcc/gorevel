package controllers

import (
	"github.com/robfig/revel"
)

type App struct {
	Application
}

func (c App) Index() revel.Result {
	return c.Render()
}
