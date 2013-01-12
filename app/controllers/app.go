package controllers

import (
	"github.com/robfig/revel"
	"SoulouUrl/app/models"
)

type Application struct {
	*rev.Controller
}

func (c Application) Index() rev.Result {
	urls, e := models.FindUrls()
	if e != nil {
		return c.RenderError(e)
	}
	return c.Render(urls)
}
