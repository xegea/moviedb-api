package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/moviedb/api/pkg/redis"
	"github.com/moviedb/api/pkg/server"
)

type Movie struct {
	Title         map[string]string `json:",omitempty"`
	Url           string            `json:",omitempty"`
	ContentRating string            `json:",omitempty"`
	Type          string            `json:",omitempty"`
	Description   map[string]string `json:",omitempty"`
	Genre         string            `json:",omitempty"`
	Image         string            `json:",omitempty"`
	DateCreated   int64             `json:",omitempty"`
	Director      []string          `json:",omitempty"`
	Actors        []string          `json:",omitempty"`
	Trailer       []Trailer         `json:",omitempty"`
}

type Trailer struct {
	Name         map[string]string `json:",omitempty"`
	Description  map[string]string `json:",omitempty"`
	Url          string            `json:",omitempty"`
	ThumbnailUrl string            `json:",omitempty"`
}

func GetMovieHandler(srv server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		b, err := redis.GetRedisValue(srv.Redis.RedisJSON, "netflix:en-gb:70090072")
		if err != nil {
			log.Println(err)
		}

		movie := Movie{}
		err = json.Unmarshal(b, &movie)
		if err != nil {
			log.Println(err)
		}

		srv.JSON(w, http.StatusOK, movie)
	}
}
