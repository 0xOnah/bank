package sqlc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}
func TestGetUsers(t *testing.T) {
	expetedUser := createRandomUser(t)
	user, err := testQueries.GetUser(context.Background(), expetedUser.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, expetedUser, user)
	require.Equal(t, expetedUser.Email, user.Email)
	require.Equal(t, expetedUser.FullName, user.FullName)
	require.Equal(t, expetedUser.HashedPassword, user.HashedPassword)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
}
