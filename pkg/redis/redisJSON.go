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
