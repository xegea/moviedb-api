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

func GetRedisValue(rh *rejson.Handler, key string) ([]byte, error) {
	value, err := redis.Bytes(rh.JSONGet(key, "."))
	if err != nil {
		return nil, err
	}

	return value, nil
}

func SetRedisValue(rh *rejson.Handler, key string, obj interface{}) (bool, error) {
	_, err := rh.JSONSet(key, ".", obj)
	if err != nil {
		return false, err
	}

	return true, nil
}
