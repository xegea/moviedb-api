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

	clients["idx:title:es"] = redisearch.NewClientFromPool(pool, "idx:title:es")
	clients["idx:title:us"] = redisearch.NewClientFromPool(pool, "idx:title:us")
	clients["idx:title:de"] = redisearch.NewClientFromPool(pool, "idx:title:de")
	clients["idx:title:gb"] = redisearch.NewClientFromPool(pool, "idx:title:gb")

	return clients
}

// func createIndexs() {

// 	// TODO Create Indexs for all available
// 	// parameters := []string{"ON", "JSON", "SCHEMA", "$.Title.US", "AS", "title", "TEXT"}
// 	// _, err = conn.Do("FT.CREATE", redis.Args{}.Add("idx:title:us").AddFlat(parameters)...)
// }

func Search(rs RediSearch, query string, page int, country string) ([]redisearch.Document, int, error) {

	docs, total, err := resolveClient(rs, country).Search(redisearch.NewQuery(fmt.Sprint("@title:", query, "*")).
		// SetReturnFields("title", "description", "type"). // if SetReturnFields
		Limit((page-1)*10, 10))
	if err != nil {
		return nil, 0, err
	}

	return docs, total, nil
}

func resolveClient(rs RediSearch, country string) *redisearch.Client {
	var client *redisearch.Client

	switch country {
	case "es":
		{
			client = rs.clients["idx:title:es"]
		}
	case "us":
		{
			client = rs.clients["idx:title:us"]
		}
	case "de":
		{
			client = rs.clients["idx:title:de"]
		}
	case "gb":
		{
			client = rs.clients["idx:title:gb"]
		}
	default:
		{
			client = rs.clients["idx:title:es"]
		}
	}
	return client
}
