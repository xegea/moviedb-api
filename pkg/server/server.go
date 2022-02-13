package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/moviedb/api/pkg/config"
	"github.com/moviedb/api/pkg/redis"
)

type Server struct {
	cfg    config.Config
	Redis  redis.Redis
	router *mux.Router
}

func NewServer(
	cfg config.Config,
	redis redis.Redis,
	r *mux.Router,
) Server {
	srv := Server{
		cfg:    cfg,
		Redis:  redis,
		router: r,
	}

	return srv
}

func (s Server) RegisterRoute(path string, handler func(w http.ResponseWriter, r *http.Request), methods []string) {
	s.router.HandleFunc(path, handler).Methods(methods...)
}

func (s Server) JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s Server) Run() error {
	addr := fmt.Sprintf(":%s", s.cfg.Port)
	if s.cfg.Env == "dev" {
		log.Printf("local env http://localhost:%s", s.cfg.Port)
		addr = fmt.Sprintf("localhost:%s", s.cfg.Port)
	}
	return http.ListenAndServe(
		addr,
		nil,
	)
}
