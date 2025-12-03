package myredis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/example/internal/config"
	"github.com/example/internal/i18n"
	"github.com/example/internal/service/jwt"
	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	Client   *redis.Client
	JService *jwt.JwtService
}

func NewRedisService(env *config.Environment, JService *jwt.JwtService) *RedisService {

	Client := redis.NewClient(&redis.Options{
		Addr:     "localhost:" + env.REDIS_PORT,
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	return &RedisService{
		Client:   Client,
		JService: JService,
	}
}

func (rs *RedisService) ValidateToken(ctx context.Context, token string) (*jwt.JWTClaims, error) {
	claim, err := rs.JService.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	res, err := rs.Client.Get(ctx, token).Result()
	if err != nil {
		return nil, err
	} else if res == "" {
		return nil, errors.New(i18n.KEY_AUTHORIZATION_UNEXPECTED_SIGNING_METHOD)
	}

	return claim, nil
}

type RedisValue struct {
	UserID    *string `json:"userId"`
	UserEmail *string `json:"userEmail"`
}

func getUserSetKey(ID *string) string {
	return fmt.Sprintf("user_id:%s", *ID)
}

func (rs *RedisService) SaveToken(ctx context.Context, token string, userID *string, userEmail *string) error {
	val := RedisValue{
		UserID:    userID,
		UserEmail: userEmail,
	}

	b, err := json.Marshal(val)
	if err != nil {
		return err
	}

	if err := rs.Client.Set(ctx, token, b, 30*24*time.Hour).Err(); err != nil {
		return err
	}

	userKey := getUserSetKey(userID)
	if err := rs.Client.SAdd(ctx, userKey, token).Err(); err != nil {
		return err
	}

	return nil
}

func (rs *RedisService) RemoveToken(ctx context.Context, token string) error {
	if rs == nil || rs.Client == nil {
		return errors.New(i18n.KEY_ERROR_REDIS_NOT_INITIALIZED)
	}
	err := rs.Client.Del(ctx, token).Err()

	return err
}

func (rs *RedisService) RemoveAllToken(ctx context.Context, UserID *string) error {
	if rs.Client == nil {
		return errors.New(i18n.KEY_ERROR_REDIS_NOT_INITIALIZED)
	}
	userKey := getUserSetKey(UserID)
	tokens, err := rs.Client.SMembers(ctx, userKey).Result()
	if err != nil {
		return err
	}

	if len(tokens) > 0 {
		if err := rs.Client.Del(ctx, tokens...).Err(); err != nil {
			return err
		}
	}

	return rs.Client.Del(ctx, userKey).Err()
}

func (rs *RedisService) CreateToken(c context.Context, userID *string, userEmail *string, userType *string) (*string, error) {
	token, err := rs.JService.GenerateToken(userID, userType)
	if err != nil {
		return nil, err
	}

	err = rs.SaveToken(c, token, userID, userEmail)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
