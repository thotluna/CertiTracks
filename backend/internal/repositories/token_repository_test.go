package repositories_test

import (
	"certitrack/internal/repositories"
	"certitrack/testutils/mocks"
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTokenRepository_RevokeToken(t *testing.T) {
	mockClient := new(mocks.MockRedisClient)
	repo := repositories.NewTokenRepository(mockClient)

	ctx := context.Background()
	token := "test-token"
	expiration := 5 * time.Minute

	statusCmd := &redis.StatusCmd{}
	statusCmd.SetVal("OK")

	mockClient.On("Set", ctx, "revoked:"+token, "1", expiration).Return(statusCmd)

	err := repo.RevokeToken(token, expiration)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestTokenRepository_IsTokenRevoked(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		setup    func(*mocks.MockRedisClient)
		expected bool
		err      error
	}{
		{
			name:  "revoked token",
			token: "revoked-token",
			setup: func(m *mocks.MockRedisClient) {
				intCmd := redis.NewIntCmd(context.Background())
				intCmd.SetVal(1)
				m.On("Exists", mock.Anything, mock.Anything).Return(intCmd)
			},
			expected: true,
			err:      nil,
		},
		{
			name:  "valid token",
			token: "valid-token",
			setup: func(m *mocks.MockRedisClient) {
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
			mockClient := new(mocks.MockRedisClient)
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
