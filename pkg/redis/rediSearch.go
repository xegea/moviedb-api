package redis

import (
	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"
)

type RediSearch struct {
	clients map[string]*redisearch.Client
}

func LoadRedisSearch(host string, password string) RediSearch {

	pool := &redis.Pool{Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", host, redis.DialPassword(password))
	}}

	clients := createClients(pool)

	rsClients := RediSearch{
		clients: clients,
	}
	return rsClients
}

func createClients(pool *redis.Pool) map[string]*redisearch.Client {

	var clients map[string]*redisearch.Client
	clients = make(map[string]*redisearch.Client)

	// TODO Create clients for every index dinamically retrieving the list of created indexs
	// indexs, err = redis.String(conn.Do("FT._LIST"))

	clients["idx:title:es"] = redisearch.NewClientFromPool(pool, "idx:title:es")
	clients["idx:title:us"] = redisearch.NewClientFromPool(pool, "idx:title:us")
	clients["idx:title:de"] = redisearch.NewClientFromPool(pool, "idx:title:de")
	clients["idx:title:gb"] = redisearch.NewClientFromPool(pool, "idx:title:gb")

	return clients
}

func createIndexs() {

	// TODO Create Indexs for all available
	// parameters := []string{"ON", "JSON", "SCHEMA", "$.Title.US", "AS", "title", "TEXT"}
	// _, err = conn.Do("FT.CREATE", redis.Args{}.Add("idx:title:us").AddFlat(parameters)...)
}
