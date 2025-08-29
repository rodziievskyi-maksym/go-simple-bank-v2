package db

import (
	"context"
	"testing"

	"github.com/rodziievskyi-maksym/go-simple-bank-v2/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, accountId int64) Entry {
	params := CreateEntryParams{
		AccountID: accountId,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), params)

	//main checks
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	//not zero checks
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	//equal
	require.Equal(t, params.AccountID, entry.AccountID)
	require.Equal(t, params.Amount, entry.Amount)

	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t, createRandomAccount(t).ID)
}

func TestGetEntry(t *testing.T) {
	randomAccount := createRandomAccount(t)
	randomEntry := createRandomEntry(t, randomAccount.ID)

	entry, err := testQueries.GetEntry(context.Background(), randomEntry.AccountID)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	require.Equal(t, randomEntry.Amount, entry.Amount)
	require.Equal(t, randomEntry.AccountID, entry.AccountID)
}

func TestListEntries(t *testing.T) {
	randomAccount := createRandomAccount(t)
	randomNum := util.RandomInt(0, 10)
	randomEntries := make([]Entry, randomNum)
	for i := 0; i < int(randomNum); i++ {
		randomEntries[i] = createRandomEntry(t, randomAccount.ID)
	}

	params := ListEntriesParams{
		AccountID: randomAccount.ID,
		Limit:     int32(randomNum),
	}
	entries, err := testQueries.ListEntries(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, entries)
	require.Len(t, entries, int(randomNum))

	for i := 0; i < int(randomNum); i++ {
		require.NotZero(t, entries[i].ID)
		require.NotZero(t, entries[i].CreatedAt)
		require.Equal(t, randomEntries[i].Amount, entries[i].Amount)
		require.Equal(t, randomEntries[i].AccountID, entries[i].AccountID)
	}
}
