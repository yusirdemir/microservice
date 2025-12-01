//go:build fiber

package server

import (
	fiberserver "github.com/yusirdemir/microservice/internal/transport/http/fiber/server"
	"github.com/yusirdemir/microservice/pkg/config"
	"go.uber.org/zap"
)

func New(cfg *config.Config, logger *zap.Logger) (Server, error) {
	return fiberserver.New(cfg, logger)
}
