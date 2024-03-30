package middleware

import (
	"arch-template/internal/app/public"
	"arch-template/pkg/errdefs"
	"arch-template/pkg/response"
	"arch-template/pkg/tlog"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

func (m *middleware) Auth(forceAuth bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := m.authLogic(ctx)
		if err != nil && (forceAuth || !isAuthFailedError(err)) {
			response.Error(ctx, err)
			return
		}
	}
}

func isAuthFailedError(err error) bool {
	return errors.Is(err, errdefs.ErrUnauthorized) || errors.Is(err, errdefs.ErrResourceNotFound)
}

func (m *middleware) authLogic(ctx *gin.Context) error {
	tokenStr := ctx.GetHeader("Authorization")
	if tokenStr == "" {
		return errdefs.ErrUnauthorized
	}
	split := strings.Split(tokenStr, " ")
	if len(split) != 2 || !strings.EqualFold(split[0], "Bearer") {
		return errdefs.ErrUnauthorized
	}
	uid, err := m.tokenManager.Verify(split[1])
	if err != nil {
		tlog.Error(ctx, "verify token failed", tlog.Fields{"err": err})
		return errdefs.ErrUnauthorized
	}
	user, err := m.userFetcher.FetchByID(ctx, uid)
	if err != nil {
		return err
	}

	ctx.Set(public.CtxAuthIDKey, uid)
	ctx.Set(public.CtxAuthModelKey, user)
	return nil
}
