package helpers

import (
	"errors"
	"fmt"
	rev "github.com/robfig/revel"
)

func getConfig(section string) (*rev.MergedConfig, error) {
  conf, e := rev.LoadConfig("app.conf")
	conf.SetSection("dev")
  if e != nil {
    return nil, errors.New("Failed to load configuration")
  }
	return conf, nil
}

func GetConfValue(key string) (string, error) {
	conf, e := getConfig("dev")
	if e != nil {
		return "", e
	}
	value, found := conf.String(key)
  if !found {
    return "", fmt.Errorf("Please set \"%s\" in \"conf/app.conf\"", key)
  }
	return value, nil
}
