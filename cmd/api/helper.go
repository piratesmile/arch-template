package main

import (
	"arch-template/configs"
	"arch-template/ent"
	"arch-template/internal/app/middleware"
	"arch-template/internal/app/user"
	"arch-template/pkg/auth"
	"arch-template/pkg/tlog"

	"go.uber.org/zap"
)

func logOptions(conf configs.Log) *tlog.Options {
	return &tlog.Options{
		Level:    conf.Level,
		LogFile:  conf.LogFile,
		FileSize: conf.FileSize,
		FileAge:  conf.FileAge,
		ErrFile:  conf.ErrFile,
	}
}

func newTokenManager(config *configs.Config) auth.TokenManager {
	return auth.NewJWT(config.Auth.SecretKey, config.Auth.Expiration)
}

func newUserFetcher(repo *user.Store) middleware.UserFetcher {
	return repo
}

func newLogger(config *configs.Config) *zap.SugaredLogger {
	return tlog.NewZapLogger(logOptions(config.Log))
}

func newDB(config *configs.Config) *ent.Client {
	orm, err := ent.Open(config.Database.Driver, config.Database.Source)
	if err != nil {
		panic(err)
	}
	if isDev(config.APP.Env) {
		orm = orm.Debug()
	}
	return orm
}

func isDev(env string) bool {
	return env == configs.Dev
}

func ginMode(env string) string {
	if isDev(env) {
		return "debug"
	} else {
		return "release"
	}
}
