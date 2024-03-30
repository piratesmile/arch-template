package app

import (
	"arch-template/internal/app/middleware"
	"arch-template/internal/app/user"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	mw middleware.Provider,
	userHandler *user.Handler,
) http.Handler {
	r := gin.Default()
	r.Use(cors.Default(), mw.RequestID(), mw.DebugLog())
	pprof.Register(r, "/debug")

	userHandler.RegisterRouter(r, mw)

	return r
}
