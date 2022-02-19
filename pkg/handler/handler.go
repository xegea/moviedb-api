package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
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

		params := mux.Vars(r)
		objid := params["movie"]
		fmt.Println("Query string key value", objid)

		b, err := redis.GetRedisValue(srv.Redis.RedisJSON, objid)
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

func SetMovieHandler(srv server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		movie := Movie{}
		if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
			srv.JSON(w, http.StatusBadRequest, "request is invalid")
			return
		}

		fmt.Println(movie.Title["GB"])

		redisKey := buildNetflixRedisKey(movie.Url)

		res, err := redis.SetRedisValue(srv.Redis.RedisJSON, redisKey, movie)
		if err != nil {
			log.Println(err)
		}

		srv.JSON(w, http.StatusOK, res)
	}
}

func buildNetflixRedisKey(movieUrl string) string {

	var redisKey string
	redisKey = strings.Replace(movieUrl, "https://www.netflix.com", "netflix", 1)
	redisKey = strings.Replace(redisKey, "/es/title/", ":es-es:", 1)
	redisKey = strings.Replace(redisKey, "/es-es/title/", ":es-es:", 1)
	redisKey = strings.Replace(redisKey, "/en-us/title/", ":en-us:", 1)
	redisKey = strings.Replace(redisKey, "/de-de/title/", ":de-de:", 1)
	redisKey = strings.Replace(redisKey, "/de/title/", ":de-de:", 1)
	redisKey = strings.Replace(redisKey, "/gb/title/", ":en-gb:", 1)
	redisKey = strings.Replace(redisKey, "/title/", ":en-us:", 1)

	if strings.Contains(redisKey, "es-en") {
		log.Fatalf("Incorrect key %s", redisKey)
	}

	if strings.Contains(redisKey, "de-en") {
		log.Fatalf("Incorrect key %s", redisKey)
	}

	return redisKey
}
