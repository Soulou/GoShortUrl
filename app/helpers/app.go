package helpers

import (
	"fmt"
)

func GetBaseUrl() (string, error) {
	addr, e := GetConfValue("http.addr")
	if e != nil {
		return "", e
	}
	port, e := GetConfValue("http.port")
	if e != nil {
		return "", e
	}
	if port == "80" {
		port = ""
	}
	return fmt.Sprintf("%s:%s", addr, port), nil
}
