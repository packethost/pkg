package testenv

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestClearAndRestore(t *testing.T) {
	nb := make([]byte, 32)
	_, err := rand.Read(nb)
	assert.NoError(t, err)
	n := base64.RawStdEncoding.EncodeToString(nb)

	vb := make([]byte, 32)
	_, err = rand.Read(vb)
	assert.NoError(t, err)
	v := base64.RawStdEncoding.EncodeToString(vb)

	assert.Empty(t, os.Getenv(n))
	os.Setenv(n, v)
	assert.Equal(t, os.Getenv(n), v)

	r := Clear()
	// ensure Clear clears the environment
	assert.Empty(t, os.Environ())
	os.Setenv(n+"AfterCleared", v)
	r.Restore()
	// ensure Restores the environment
	assert.NotEmpty(t, os.Environ())
	assert.Equal(t, os.Getenv(n), v)
	// and that an env set after Clear does not persist after Restore is called
	assert.Empty(t, os.Getenv(n+"AfterCleared"))
}
