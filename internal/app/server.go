package app

import (
	"arch-template/configs"
	"arch-template/ent"
	"context"
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	conf      *configs.Config
	orm       *ent.Client
	apiServer *http.Server
	logger    *zap.SugaredLogger
}

func NewServer(config *configs.Config, logger *zap.SugaredLogger, orm *ent.Client, handler http.Handler) *Server {
	return &Server{
		conf:   config,
		logger: logger,
		orm:    orm,
		apiServer: &http.Server{
			Addr:    config.APP.Host,
			Handler: handler,
		},
	}
}

func (s *Server) Run() {
	s.logger.Infow("Start to listening http server", "addr", s.apiServer.Addr)
	go func() {
		if err := s.apiServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Panicw("HTTP server listen failed", "err", err)
		}
	}()
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := s.apiServer.Shutdown(ctx); err != nil {
		s.logger.Errorw("shutdown api server failed", "err", err)
	}

	if err := s.orm.Close(); err != nil {
		s.logger.Errorw("close db orm failed", "err", err)
	}
}
