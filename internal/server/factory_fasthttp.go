//go:build fasthttp

package server

import (
	fasthttpserver "github.com/yusirdemir/microservice/internal/transport/http/fasthttp/server"
	"github.com/yusirdemir/microservice/pkg/config"
	"go.uber.org/zap"
)

func New(cfg *config.Config, logger *zap.Logger) (Server, error) {
	return fasthttpserver.New(cfg, logger)
}
