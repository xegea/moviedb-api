package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/moviedb/api/pkg/redis"
	"github.com/moviedb/api/pkg/server"
)

type Movie struct {
	Title         string    `json:",omitempty"`
	Url           string    `json:",omitempty"`
	ContentRating string    `json:",omitempty"`
	Type          string    `json:",omitempty"`
	Description   string    `json:",omitempty"`
	Genre         string    `json:",omitempty"`
	Image         string    `json:",omitempty"`
	ReleaseDate   int64     `json:",omitempty"`
	Director      []string  `json:",omitempty"`
	Actors        []string  `json:",omitempty"`
	Trailer       []Trailer `json:",omitempty"`
	Updated       int64     `json:",omitempty"`
}

type Trailer struct {
	Name         string `json:",omitempty"`
	Description  string `json:",omitempty"`
	Url          string `json:",omitempty"`
	ThumbnailUrl string `json:",omitempty"`
}

func SetMovieHandler(srv server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("ApiKey") != srv.Config.AdminApiKey {
			srv.JSON(w, http.StatusForbidden, errors.New("wrong api key"))
			return
		}

		movie := Movie{}
		if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
			srv.JSON(w, http.StatusBadRequest, err)
			return
		}

		redisKey, err := buildNetflixRedisKey(movie.Url)
		if err != nil {
			srv.JSON(w, http.StatusInternalServerError, "")
			return
		}

		_, err = redis.SetRedisValue(srv.Redis.RedisJSON, redisKey, movie)
		if err != nil {
			srv.JSON(w, http.StatusInternalServerError, err)
			return
		}

		srv.JSON(w, http.StatusCreated, nil)
	}
}

func SearchHandler(srv server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("ApiKey") != srv.Config.ApiKey && r.Header.Get("ApiKey") != srv.Config.AdminApiKey {
			srv.JSON(w, http.StatusForbidden, errors.New("wrong api key"))
			return
		}

		var page int
		// TODO rowsPerPage by config
		rowsPerPage := 8

		fmt.Sscan(r.URL.Query().Get("page"), &page)
		query := r.URL.Query().Get("query")
		culture := r.URL.Query().Get("culture")

		docs, total, err := redis.Search(srv.Redis.RediSearch, query, page, rowsPerPage, culture)
		if err != nil {
			srv.JSON(w, http.StatusInternalServerError, err)
			return
		}

		var movieList []Movie

		for _, v := range docs {

			var jsonbody string
			for _, p := range v.Properties {
				jsonbody = fmt.Sprint(p)
			}
			if err != nil {
				log.Fatal(err)
			}

			var content Movie
			if err := json.Unmarshal([]byte(jsonbody), &content); err != nil {
				log.Fatalf("Failed to Unmarshall %s", jsonbody)
			}
			movieList = append(movieList, content)
		}

		type PagedMovieList struct {
			MovieList []Movie
			LastPage  bool
		}

		var lastPage = IsLastPage(page, rowsPerPage, total)
		var pagedMovieList PagedMovieList = PagedMovieList{
			MovieList: movieList,
			LastPage:  lastPage,
		}

		srv.JSON(w, http.StatusOK, pagedMovieList)
	}
}

func IsLastPage(page, rowsPerPage, total int) bool {
	totalPages := math.Ceil((float64(total)/float64(rowsPerPage))*1) / 1

	return float64(page) >= totalPages
}

func GetMovieHandler(srv server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("ApiKey") != srv.Config.ApiKey && r.Header.Get("ApiKey") != srv.Config.AdminApiKey {
			srv.JSON(w, http.StatusForbidden, errors.New("wrong api key"))
			return
		}

		id, found := mux.Vars(r)["id"]
		if !found {
			srv.JSON(w, http.StatusBadRequest, nil)
			return
		}

		b, err := redis.GetRedisValue(srv.Redis.RedisJSON, id)
		if err != nil {
			srv.JSON(w, http.StatusNotFound, err)
			return
		}

		movie := Movie{}
		err = json.Unmarshal(b, &movie)
		if err != nil {
			srv.JSON(w, http.StatusInternalServerError, err)
			return
		}

		srv.JSON(w, http.StatusOK, movie)
	}
}

func buildNetflixRedisKey(movieUrl string) (string, error) {

	var key string
	key = strings.Replace(movieUrl, "https://www.netflix.com", "netflix", 1)
	key = strings.Replace(key, "/es/title/", ":es-es:", 1)
	key = strings.Replace(key, "/es-es/title/", ":es-es:", 1)
	key = strings.Replace(key, "/es-en/title/", ":es-en:", 1)
	key = strings.Replace(key, "/en-us/title/", ":en-us:", 1)
	key = strings.Replace(key, "/de-de/title/", ":de-de:", 1)
	key = strings.Replace(key, "/de/title/", ":de-de:", 1)
	key = strings.Replace(key, "/gb/title/", ":en-gb:", 1)
	key = strings.Replace(key, "/title/", ":en-us:", 1)

	if strings.Contains(key, "de-en") {
		return "", errors.New("incorrect key " + key)
	}

	return key, nil
}
