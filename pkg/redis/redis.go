package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson"
)

type Redis struct {
	Conn       redis.Conn
	RedisJSON  *rejson.Handler
	RediSearch RediSearch
}

func Connect(host string, password string) (Redis, error) {

	conn, err := redis.Dial("tcp", host, redis.DialPassword(password))
	if err != nil {
		return Redis{}, err
	}

	redisJSON := LoadRedisJSON(conn)
	redisSearch := LoadRedisSearch(host, password)

	redisConn := Redis{
		Conn:       conn,
		RedisJSON:  redisJSON,
		RediSearch: redisSearch,
	}

	return redisConn, nil
}

func Close(conn redis.Conn) {
	conn.Close()
}
