//go:build std

package server

import (
	stdserver "github.com/yusirdemir/microservice/internal/transport/http/std/server"
	"github.com/yusirdemir/microservice/pkg/config"
	"go.uber.org/zap"
)

func New(cfg *config.Config, logger *zap.Logger) (Server, error) {
	return stdserver.New(cfg, logger)
}
