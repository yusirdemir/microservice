//go:build echo

package server

import (
	echoserver "github.com/yusirdemir/microservice/internal/transport/http/echo/server"
	"github.com/yusirdemir/microservice/pkg/config"
	"go.uber.org/zap"
)

func New(cfg *config.Config, logger *zap.Logger) (Server, error) {
	return echoserver.New(cfg, logger)
}
