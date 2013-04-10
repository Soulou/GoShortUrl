package helpers

import (
	"errors"
	redis "github.com/alphazero/Go-Redis"
	"strconv"
)

func GetRedisClient() (redis.Client, error) {
	redis_pass, e := GetConfValue("redis.password")
	if e != nil {
		redis_pass = ""
	}
	redis_db, e := GetConfValue("redis.dbnumber")
	if e != nil {
		return nil, e
	}

	redis_ndb, _ := strconv.ParseInt(redis_db, 10, 32)
  spec := redis.DefaultSpec().Db(int(redis_ndb)).Password(redis_pass)
  client, e := redis.NewSynchClientWithSpec(spec)
  if e != nil {
    return nil, errors.New("Failed to create redis client, wrong password ?")
  }
	return client, nil
}
