// Package repositories provides a centralized way to register all custom validators.
package repositories

import (
	"context"
	"time"
)

type TokenRepository interface {
	RevokeToken(token string, expiresIn time.Duration) error
	IsTokenRevoked(token string) (bool, error)
}

type tokenRepository struct {
	client RedisClient
}

func NewTokenRepository(client RedisClient) TokenRepository {
	return &tokenRepository{client: client}
}

func (r *tokenRepository) RevokeToken(token string, expiresIn time.Duration) error {
	key := "revoked:" + token
	return r.client.Set(context.Background(), key, "1", expiresIn).Err()
}

func (r *tokenRepository) IsTokenRevoked(token string) (bool, error) {
	key := "revoked:" + token
	result, err := r.client.Exists(context.Background(), key).Result()
	return result > 0, err
}
