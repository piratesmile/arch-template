package middleware

import (
	"arch-template/configs"
	"arch-template/internal/app/entity"
	"arch-template/pkg/auth"
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type Provider interface {
	Auth(force bool) gin.HandlerFunc
	RequestID() gin.HandlerFunc
	Limit(rps rate.Limit, count int, ttl time.Duration) gin.HandlerFunc
	PopulateLogger() gin.HandlerFunc
}

type UserFetcher interface {
	FetchByID(ctx context.Context, id uint) (*entity.User, error)
}

type middleware struct {
	config       *configs.Config
	tokenManager auth.TokenManager
	userFetcher  UserFetcher
	logger       *zap.SugaredLogger
}

func New(config *configs.Config, tokenManager auth.TokenManager, userFetcher UserFetcher, logger *zap.SugaredLogger) Provider {
	return &middleware{
		config:       config,
		tokenManager: tokenManager,
		userFetcher:  userFetcher,
		logger:       logger,
	}
}
