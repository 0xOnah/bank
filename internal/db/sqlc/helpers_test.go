package sqlc

import (
	"context"
	"fmt"
	"testing"

	"github.com/0xOnah/bank/internal/sdk/auth"
	"github.com/0xOnah/bank/internal/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) *User {
	username := util.RandomOwner()
	password, err := auth.HashPassword(util.RandomString(10))
	require.NoError(t, err)
	user, err := testQueries.CreateUser(context.Background(), CreateUserParams{
		Username:       username,
		HashedPassword: password,
		FullName:       fmt.Sprintf("%s %s", util.RandomOwner(), util.RandomOwner()),
		Email:          fmt.Sprintf("%s@gmail.com", username),
	})
	require.NoError(t, err)
	return user
}

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return *account
}

func createRandomEntry(t *testing.T) Entry {
	account := createRandomAccount(t)
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotEmpty(t, entry.ID)
	require.NotEmpty(t, entry.CreatedAt)

	return *entry
}
