package app

import (
	"arch-template/configs"
	"arch-template/ent"
	"arch-template/pkg/tlog"
	"context"
	"errors"
	"net/http"
	"time"
)

type Server struct {
	conf      *configs.Config
	orm       *ent.Client
	apiServer *http.Server
}

func NewServer(config *configs.Config, orm *ent.Client, handler http.Handler) *Server {
	return &Server{
		conf: config,
		orm:  orm,
		apiServer: &http.Server{
			Addr:    config.APP.Host,
			Handler: handler,
		},
	}
}

func (s *Server) Run() {
	tlog.Info(context.Background(), "Start to listening http server", tlog.Fields{"addr": s.apiServer.Addr})
	go func() {
		if err := s.apiServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			tlog.Panic(context.Background(), "HTTP server listen failed", tlog.Fields{"err": err})
		}
	}()
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := s.apiServer.Shutdown(ctx); err != nil {
		tlog.Error(ctx, "shutdown api server failed", tlog.Fields{"err": err})
	}

	if err := s.orm.Close(); err != nil {
		tlog.Error(ctx, "close db orm failed", tlog.Fields{"err": err})
	}
}
