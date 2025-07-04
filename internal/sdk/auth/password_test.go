package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const password = "secret"

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword(password)
	require.NoError(t, err)
	require.NotZero(t, hash)

	//correct password
	ok := ComparePassword([]byte(hash), password)
	require.True(t, ok)

	//compare wrong password
	ok = ComparePassword([]byte("hello"), password)
	require.False(t, ok)

	//checking to ensure we don't have the same hashes
	hash1, err := HashPassword(password)
	require.NoError(t, err)

	require.NotEqual(t, hash, hash1)
}
