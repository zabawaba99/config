package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestParseJSON(t *testing.T) {
	filename := "example/config.json"
	config, err := parseJSON(filename)
	require.NoError(t, err)

	assert.Len(t, config, 2)

	arg1, exists := config["port"]
	require.True(t, exists)
	// JSON marshals all numbers to float64
	assert.EqualValues(t, float64(8080), arg1.Default)
	assert.IsType(t, float64(0), arg1.Default)
	assert.Equal(t, "PORT", arg1.EnvName)
	assert.Equal(t, "port", arg1.FlagName)
	assert.Equal(t, "The port that the application will run on", arg1.Description)

	arg2, exists := config["s3_bucket"]
	require.True(t, exists)
	assert.Equal(t, "", arg2.Default)
	assert.Equal(t, "S3_BUCKET", arg2.EnvName)
	assert.Equal(t, "s3.bucket", arg2.FlagName)
	assert.Equal(t, "The s3 bucket used to upload icons", arg2.Description)
}

func TestParseJSONFailure(t *testing.T) {
	_, err := parseJSON("foobar")
	require.Error(t, err)
}
