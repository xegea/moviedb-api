package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson"
)

type Redis struct {
	Conn      redis.Conn
	redisJSON *rejson.Handler
	//rediSearch redisSearchClient
}

func Connect(host string, password string) (Redis, error) {

	conn, err := redis.Dial("tcp", host, redis.DialPassword(password))
	if err != nil {
		return Redis{}, err
	}

	redisJSON := LoadRedisJSON(conn)

	//redisSearchClient := LoadRedisSearch(host, password)

	redisConn := Redis{
		Conn:      conn,
		redisJSON: redisJSON,
		//rediSearch: redisSearchClient,
	}

	return redisConn, nil
}

func Close(conn redis.Conn) {
	conn.Close()
}
