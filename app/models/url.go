package models

import (
	revel "github.com/robfig/revel"
	helpers "github.com/Soulou/GoShortUrl/app/helpers"
	"crypto/sha1"
	"encoding/binary"
	"io"
	"fmt"
	"html"
	"html/template"
)

type Url struct {
	Href string
	Digest string
	Counter uint64
}

func init() {
	revel.TemplateFuncs["url_string"] = func(url Url) template.HTML {
		addr, e := helpers.GetBaseUrl()
		if e != nil {
			return template.HTML("Error loading base URL")
		}
		href := url.Href
		if len(href) > 64 {
			begin_url := href[:30]
			end_url := href[len(href)-30:len(href)-1]
			href = fmt.Sprintf("%s[...]%s", begin_url, end_url)
		}
		return template.HTML(fmt.Sprintf("<a href=\"%s\">%s</a> : <a href=\"http://%s/%s\">http://%s/%s</a>",
			url.Href, html.EscapeString(href), addr, url.Digest, addr, url.Digest))
	}
}

func NewUrl(url string) Url {
	h := sha1.New()
	io.WriteString(h, url)

	res := Url{url, fmt.Sprintf("%x", h.Sum(nil)), 0}

	return res
}

func (url Url) IncCounter() error {
	redis_client, e := helpers.GetRedisClient()
	if e != nil {
		return e
	}
	redis_client.Incr(fmt.Sprintf("visits:%s", url.Digest))
  return nil
}

func FindUrl(digest string) (*Url, error) {
	redis_client, e := helpers.GetRedisClient()
  if e != nil {
		return nil, e
	}

  // Parrallelize both get
  c_href := make(chan []byte)
  go func() {
    href, _ := redis_client.Get(digest)
    c_href <- href
  }()
  c_visits := make(chan []byte)
  go func() {
    visits, _ := redis_client.Get(fmt.Sprintf("visits:%s", digest))
    c_visits <- visits
  }()

  i_visits, _ := binary.Uvarint(<-c_visits)
	url := Url{string(<-c_href), digest, i_visits}

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

    url, e_url := FindUrl(string(digest))
		if e_url != nil {
			return nil, e_url
		}
    urls[i] = *url
	}
	return urls[:], nil
}

func (u Url) Save() error {
	redis_client, e := helpers.GetRedisClient()
	if e != nil {
		return e
	}

  buf := make([]byte, 8)
  binary.PutUvarint(buf, uint64(0))
	redis_client.Set(u.Digest[0:5], []byte(u.Href))
  redis_client.Set(fmt.Sprintf("visits:%s", u.Digest[0:5]), buf)
	redis_client.Rpush("digests", []byte(u.Digest[0:5]))

	return nil
}
