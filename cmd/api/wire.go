//go:build wireinject
// +build wireinject

package main

import (
	"arch-template/configs"
	"arch-template/internal/app"
	"arch-template/internal/app/middleware"
	"arch-template/internal/app/public"
	"arch-template/internal/app/user"

	"github.com/google/wire"
)

var baseSet = wire.NewSet(
	public.NewBaseRepository,
)

var moduleSet = wire.NewSet(
	user.NewHandler, user.NewService, user.NewRepository,
)

func newServer(c *configs.Config) *app.Server {
	panic(wire.Build(
		baseSet,
		moduleSet,
		newUserFetcher,
		newTokenManager,
		newDB,
		middleware.New,
		app.NewRouter,
		app.NewServer,
	))
}
