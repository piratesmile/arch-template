package user

import (
	"arch-template/internal/app/middleware"
	"arch-template/internal/app/public"
	"arch-template/pkg/errdefs"
	"arch-template/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (u *Handler) RegisterRouter(router gin.IRouter, mw middleware.Provider) {
	router.POST("login", u.Login)

	authRouter := router.Group("", mw.Auth(true))
	{
		authRouter.GET("profile", u.Profile)
	}
}

func (u *Handler) Login(ctx *gin.Context) {
	var req LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, errdefs.InvalidParams(err))
		return
	}
	resp, err := u.service.Login(ctx, req)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, resp)
}

func (u *Handler) Profile(ctx *gin.Context) {
	auth, err := public.AuthModelFromContext(ctx)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, auth)
}
