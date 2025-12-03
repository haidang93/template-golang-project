package myredis

import (
	"context"

	"github.com/example/internal/service/jwt"
)

type RedisServiceInterface interface {
	ValidateToken(c context.Context, token string) (*jwt.JWTClaims, error)
	RemoveToken(c context.Context, token string) error
	RemoveAllToken(ctx context.Context, UserID *string) error
	CreateToken(c context.Context, userID *string, userEmail *string, userType *string) (*string, error)
}
