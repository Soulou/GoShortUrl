package controllers

import (
	"github.com/robfig/revel"
	"github.com/Soulou/GoShortUrl/app/models"
)

type UrlController struct {
  *rev.Controller
}

func (c UrlController) Save(url string) rev.Result {
	u := models.NewUrl(url)
	e := u.Save()
	if e != nil {
		return c.RenderError(e)
	}
  return c.Redirect("/")
}

func (c UrlController) Show(digest string) rev.Result {
	url,e := models.FindUrl(digest)
	if e != nil {
		return c.RenderError(e)
	}
	if(c.Request.Format == "json") {
		return c.RenderJson(url)
	}
	return c.Redirect(url.Href)
}
