package redis

import (
	"fmt"

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

	redisSearch := RediSearch{
		clients: clients,
	}
	return redisSearch
}

func createClients(pool *redis.Pool) map[string]*redisearch.Client {

	clients := make(map[string]*redisearch.Client)

	// TODO Create clients for every index dinamically retrieving the list of created indexs
	// indexs, err = redis.String(conn.Do("FT._LIST"))

	clients["idx:title:es-es"] = redisearch.NewClientFromPool(pool, "idx:title:es-es")
	clients["idx:title:en-es"] = redisearch.NewClientFromPool(pool, "idx:title:en-es")
	clients["idx:title:en-us"] = redisearch.NewClientFromPool(pool, "idx:title:en-us")
	clients["idx:title:de-de"] = redisearch.NewClientFromPool(pool, "idx:title:de-de")
	clients["idx:title:en-gb"] = redisearch.NewClientFromPool(pool, "idx:title:en-gb")

	return clients
}

// func createIndexs() {

// 	// TODO Create Indexs for all available
// 	// parameters := []string{"ON", "JSON", "SCHEMA", "$.Title.US", "AS", "title", "TEXT"}
// 	// _, err = conn.Do("FT.CREATE", redis.Args{}.Add("idx:title:us").AddFlat(parameters)...)
// }

func Search(rs RediSearch, query string, page int, culture string) ([]redisearch.Document, int, error) {

	docs, total, err := resolveClient(rs, culture).Search(redisearch.NewQuery(fmt.Sprint("@title:", query, "*")).
		// SetReturnFields("title", "description", "type"). // if SetReturnFields
		Limit((page-1)*10, 10))
	if err != nil {
		return nil, 0, err
	}

	return docs, total, nil
}

func resolveClient(rs RediSearch, culture string) *redisearch.Client {
	var client *redisearch.Client

	switch culture {
	case "es-es", "en-es", "en-us", "de-de", "en-gb":
		{
			client = rs.clients["idx:title:"+culture]
		}
	default:
		{
			client = rs.clients["idx:title:es-es"]
		}
	}
	return client
}
