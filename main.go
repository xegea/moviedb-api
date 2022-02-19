package main

import (
	"flag"
	"log"

	"github.com/gorilla/mux"
	"github.com/moviedb/api/pkg/config"
	"github.com/moviedb/api/pkg/handler"
	"github.com/moviedb/api/pkg/redis"
	"github.com/moviedb/api/pkg/server"
)

func main() {

	env := flag.String("env", ".env", "environment path")
	flag.Parse()

	cfg, err := config.LoadConfig(env)
	if err != nil {
		log.Fatalf("unable to load config: %+v", err)
	}

	redisConn, err := redis.Connect(cfg.RedisHost, cfg.RedisPassword)
	if err != nil {
		log.Fatalf("unable to connect to redis: %v", err)
	}
	defer redis.Close(redisConn.Conn)

	srv := server.NewServer(
		cfg,
		redisConn,
		mux.NewRouter(),
	)

	srv.RegisterRoute("/movie/{id}", handler.GetMovieHandler(srv), []string{"GET"})
	srv.RegisterRoute("/movie", handler.SetMovieHandler(srv), []string{"POST"})
	srv.RegisterRoute("/search/", handler.SearchHandler(srv), []string{"GET"})

	log.Fatal(srv.Run())
}
