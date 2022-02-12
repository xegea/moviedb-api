package server

import (
	"github.com/moviedb/api/pkg/config"
	"github.com/moviedb/api/pkg/redis"

	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg    config.Config
	redis  redis.Redis
	router *gin.Engine
}

func NewServer(
	cfg config.Config,
	redis redis.Redis,
	router *gin.Engine,
) Server {
	srv := Server{
		cfg:    cfg,
		redis:  redis,
		router: router,
	}
	return srv
}
