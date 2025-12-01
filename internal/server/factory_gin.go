//go:build gin

package server

import (
	ginserver "github.com/yusirdemir/microservice/internal/transport/http/gin/server"
	"github.com/yusirdemir/microservice/pkg/config"
	"go.uber.org/zap"
)

func New(cfg *config.Config, logger *zap.Logger) (Server, error) {
	return ginserver.New(cfg, logger)
}
