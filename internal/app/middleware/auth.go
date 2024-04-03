package middleware

import (
	"arch-template/internal/app/entity"
	"arch-template/pkg/errdefs"
	"arch-template/pkg/response"
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const contextKeyUser = "ctx-user"

func (m *middleware) Auth(forceAuth bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := m.realAuth(ctx)
		if err != nil && (forceAuth || !isAuthFailedError(err)) {
			response.Error(ctx, err)
			return
		}
	}
}

func isAuthFailedError(err error) bool {
	return errors.Is(err, errdefs.ErrUnauthorized) || errors.Is(err, errdefs.ErrResourceNotFound)
}

func (m *middleware) realAuth(ctx *gin.Context) error {
	token := tokenFromHeader(ctx.Request)
	if token == "" {
		return errdefs.ErrUnauthorized
	}
	uid, err := m.tokenManager.Verify(token)
	if err != nil {
		return errdefs.ErrUnauthorized
	}
	user, err := m.userFetcher.FetchByID(ctx, uid)
	if err != nil {
		return err
	}

	ctx.Set(contextKeyUser, user)
	return nil
}

func tokenFromHeader(r *http.Request) string {
	header := r.Header.Get("Authorization")
	if len(header) > 7 && strings.ToLower(header[:6]) == "bearer" {
		return header[:7]
	}
	return ""
}

func UserFromContext(ctx context.Context) (*entity.User, error) {
	m := ctx.Value(contextKeyUser)
	if v, ok := m.(*entity.User); ok && v.ID > 0 {
		return v, nil
	}
	return nil, errdefs.ErrUnauthorized
}
