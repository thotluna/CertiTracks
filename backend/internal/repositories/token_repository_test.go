package repositories_test

import (
	"certitrack/internal/repositories"
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
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
	args := m.Called(ctx, keys)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestTokenRepository_RevokeToken(t *testing.T) {
	mockClient := new(MockRedisClient)
	repo := repositories.NewTokenRepository(mockClient)

	ctx := context.Background()
	token := "test-token"
	expiration := 5 * time.Minute

	// Crear un StatusCmd mock
	statusCmd := &redis.StatusCmd{}
	statusCmd.SetVal("OK")

	// Configurar las expectativas
	mockClient.On("Set", ctx, "revoked:"+token, "1", expiration).Return(statusCmd)

	err := repo.RevokeToken(token, expiration)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestTokenRepository_IsTokenRevoked(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		setup    func(*MockRedisClient)
		expected bool
		err      error
	}{
		{
			name:  "token revocado",
			token: "revoked-token",
			setup: func(m *MockRedisClient) {
				intCmd := redis.NewIntCmd(context.Background())
				intCmd.SetVal(1)
				m.On("Exists", mock.Anything, mock.Anything).Return(intCmd)
			},
			expected: true,
			err:      nil,
		},
		{
			name:  "token no revocado",
			token: "valid-token",
			setup: func(m *MockRedisClient) {
				intCmd := redis.NewIntCmd(context.Background())
				intCmd.SetVal(0)
				m.On("Exists", mock.Anything, mock.Anything).Return(intCmd)
			},
			expected: false,
			err:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockRedisClient)
			repo := repositories.NewTokenRepository(mockClient)

			if tt.setup != nil {
				tt.setup(mockClient)
			}

			revoked, err := repo.IsTokenRevoked(tt.token)

			if tt.err != nil {
				assert.ErrorIs(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, revoked)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
