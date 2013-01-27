package models

import (
	revel "github.com/robfig/revel"
	helpers "github.com/Soulou/GoShortUrl/app/helpers"
	"crypto/sha1"
	"io"
	"fmt"
	"html"
	"html/template"
)

type Url struct {
	Href string
	Digest string
	Counter int
}

func init() {
	revel.TemplateFuncs["url_string"] = func(url Url) template.HTML {
		addr, e := helpers.GetBaseUrl()
		if e != nil {
			return template.HTML("Error loading base URL")
		}
		href := url.Href
		if len(href) > 64 {
			href = fmt.Sprintf("%s[...]%s", href[:30], href[len(href)-30:len(href)-1])
		}
		return template.HTML(fmt.Sprintf("<a href=\"%s\">%s</a> : <a href=\"http://%s/%x\">http://%s/%x</a>",
			url.Href, html.EscapeString(href), addr, url.Digest, addr, url.Digest))
	}
}

func NewUrl(url string) Url {
	h := sha1.New()
	io.WriteString(h, url)

	res := Url{url, fmt.Sprintf("%x", h.Sum(nil)), 0}

	return res
}

func (url Url) IncCounter() {
	redis_client, e := helpers.GetRedisClient()
	if e != nil {
		return nil, e
	}
	redis_client.Inc(fmt.Sprintf("visits:%s", url.Digest))
}

func FindUrl(digest string) (*Url, error) {
	redis_client, e := helpers.GetRedisClient()
	if e != nil {
		return nil, e
	}

	href, e := redis_client.Get(digest)
	if e != nil {
		return nil , e
	}
	url := Url{string(href), digest}

	return &url, nil
}

func FindUrls() ([]Url, error) {
	redis_client, e := helpers.GetRedisClient()
	if e != nil {
		return nil, e
	}

	nb_urls, e := redis_client.Llen("digests")
	if e!= nil {
		return nil, e
	}

	urls := make([]Url, nb_urls)
	for i := int64(0); i < nb_urls; i++ {
		digest, e := redis_client.Lindex("digests", i)
		if e != nil {
			return nil, e
		}

		url_bytes, e := redis_client.Get(fmt.Sprintf("%x", digest))
		url := fmt.Sprintf("%s", url_bytes)
		urls[i] = Url{url,fmt.Sprintf("%x", digest)}
	}
	return urls[:], nil
}

func (u Url) Save() error {
	redis_client, e := helpers.GetRedisClient()
	if e != nil {
		return e
	}
	redis_client.Set(u.Digest[0:5], []byte(u.Href))
	redis_client.Set(fmt.Sprintf("%s_cnt", u.Digest[0:5]), 0)
	redis_client.Rpush("digests", []byte(u.Digest[0:5]))

	return nil
}
