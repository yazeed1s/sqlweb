package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgs_ValidatePortRange_ValidPort(t *testing.T) {
	args := &Args{
		Port: 3000,
	}

	err := args.ValidatePortRange()

	assert.NoError(t, err, "Expected no error for valid port")
}

func TestArgs_ValidatePortRange_InvalidPort(t *testing.T) {
	args := &Args{
		Port: -1,
	}

	err := args.ValidatePortRange()

	assert.Error(t, err, "Expected error for invalid port")
	assert.EqualError(t, err, "invalid port number", "Expected error message for invalid port")
}

func TestArgs_NewArgs_DefaultValues(t *testing.T) {
	args := NewArgs()

	assert.Equal(t, 3000, args.Port, "Expected default port value to be 3000")
	assert.True(t, args.Log, "Expected default log value to be true")
	assert.Contains(t, args.Help, "USAGE: sqlweb", "Expected default help message to contain usage information")
	assert.Equal(t, "version 0.1.0", args.Version, "Expected default version to be 'version 0.1.0'")
}

func TestArgs_NewArgs_SetCustomValues(t *testing.T) {
	args := NewArgs()
	args.Port = 8080
	args.Log = false
	args.Version = "1.2.3"
	assert.Equal(t, 8080, args.Port, "Expected custom port value to be set")
	assert.False(t, args.Log, "Expected custom log value to be set")
	assert.Contains(t, args.Help, "USAGE: sqlweb", "Expected default help message to contain usage information")
	assert.Equal(t, "1.2.3", args.Version, "Expected custom version to be set")
}
