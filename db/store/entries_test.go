package db

import (
	"context"
	"testing"
	"time"

	"github.com/onahvictor/bank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	arg := createRandomEntry(t)

	entry, err := testQueries.GetEntry(context.Background(), arg.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.ID, entry.ID)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.WithinDuration(t, arg.CreatedAt, entry.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	entry := createRandomEntry(t)
	//inserting more entries for the same account
	arg := CreateEntryParams{
		AccountID: entry.AccountID,
		Amount:    util.RandomMoney(),
	}

	for i := 0; i < 10; i++ { 
		value, err := testQueries.CreateEntry(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, value)
	}

	params := ListEntriesParams{
		AccountID: entry.AccountID,
		Limit:     5,
		Offset:    0,
	}
	entries, err := testQueries.ListEntries(context.Background(), params)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, account := range entries {
		require.NotEmpty(t, account)
	}
}
