package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"certitrack/internal/config"
	"certitrack/testutils"
)

func TestNewClient(t *testing.T) {
	testCfg := testutils.GetTestConfig()

	tests := []struct {
		name    string
		cfg     *config.RedisConfig
		wantErr bool
	}{
		{
			name:    "nil config",
			cfg:     nil,
			wantErr: true,
		},
		{
			name: "invalid URL",
			cfg: &config.RedisConfig{
				URL: "invalid-url",
			},
			wantErr: true,
		},
		{
			name:    "valid config",
			cfg:     &testCfg.Redis,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.cfg)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				require.NoError(t, err)
				require.NotNil(t, client)
				assert.NoError(t, client.Close())
			}
		})
	}
}

func TestClient_Close(t *testing.T) {
	// Test closing a nil client
	t.Run("nil client", func(t *testing.T) {
		var c *Client
		err := c.Close()
		assert.NoError(t, err)
	})

	// Test closing a valid client
	t.Run("valid client", func(t *testing.T) {
		cfg := testutils.GetTestConfig()
		client, err := NewClient(&cfg.Redis)
		require.NoError(t, err)
		require.NotNil(t, client)

		err = client.Close()
		assert.NoError(t, err)
	})
}

func TestClient_Client(t *testing.T) {
	cfg := testutils.GetTestConfig()
	client, err := NewClient(&cfg.Redis)
	require.NoError(t, err)
	defer client.Close()

	// Test that Client() returns a non-nil client
	assert.NotNil(t, client.Client())
}
