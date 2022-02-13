package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson"
)

func LoadRedisJSON(conn redis.Conn) *rejson.Handler {
	reJsonClient := rejson.NewReJSONHandler()
	reJsonClient.SetRedigoClient(conn)

	return reJsonClient
}

func GetRedisValue(rh *rejson.Handler, redisKey string) ([]byte, error) {
	redisValue, err := redis.Bytes(rh.JSONGet(redisKey, "."))
	if err != nil {
		return nil, err
	}

	return redisValue, nil
}
