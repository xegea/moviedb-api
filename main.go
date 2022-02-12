package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/moviedb/api/pkg/config"
	"github.com/moviedb/api/pkg/redis"
	"github.com/moviedb/api/pkg/server"

	"github.com/gin-gonic/gin"
)

var (
	env string
)

func main() {
	flag.Parse()

	cfg, err := config.LoadConfig(env)
	if err != nil {
		log.Fatalf("unable to load config: %+v", err)
	}

	r, err := redis.Connect(cfg.RedisHost, cfg.RedisPassword)
	if err != nil {
		log.Fatalf("unable to connect to redis: %v", err)
	}
	defer redis.Close(r.Conn)

	srv := server.NewServer(
		cfg,
		r,
		gin.Default(),
	)

	fmt.Println(srv)
}
