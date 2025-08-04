// Package mocks provides test doubles for external dependencies.
// It includes mock implementations of interfaces used throughout the application
// to facilitate isolated unit testing without requiring actual external services.
package mocks

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
)

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	return args.Get(0).(*redis.StatusCmd)
}

func (m *MockRedisClient) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	callArgs := m.Called(append([]interface{}{ctx}, stringSliceToInterfaceSlice(keys)...)...)
	return callArgs.Get(0).(*redis.IntCmd)
}

func stringSliceToInterfaceSlice(strSlice []string) []interface{} {
	ifaceSlice := make([]interface{}, len(strSlice))
	for i, v := range strSlice {
		ifaceSlice[i] = v
	}
	return ifaceSlice
}

func (m *MockRedisClient) Close() error {
	args := m.Called()
	return args.Error(0)
}
