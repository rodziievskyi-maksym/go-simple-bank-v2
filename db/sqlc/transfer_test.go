package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/rodziievskyi-maksym/go-simple-bank-v2/util"
	"github.com/stretchr/testify/require"
)

//create
//get
//list

func createRandomTransfer(t *testing.T, fromAccountId, toAccountId int64) Transfer {
	params := CreateTransferParams{
		FromAccountID: sql.NullInt64{Int64: fromAccountId, Valid: true},
		ToAccountID:   sql.NullInt64{Int64: toAccountId, Valid: true},
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	//non zero
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	require.Equal(t, params.FromAccountID, transfer.FromAccountID)
	require.Equal(t, params.ToAccountID, transfer.ToAccountID)
	require.Equal(t, params.Amount, transfer.Amount)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	createRandomTransfer(t, account1.ID, account2.ID)
}

func TestGetTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	randomTransfer := createRandomTransfer(t, account1.ID, account2.ID)
	transfer, err := testQueries.GetTransfer(context.Background(), randomTransfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, randomTransfer.ID, transfer.ID)
	require.Equal(t, randomTransfer.FromAccountID, transfer.FromAccountID)
	require.Equal(t, randomTransfer.ToAccountID, transfer.ToAccountID)
	require.Equal(t, randomTransfer.Amount, transfer.Amount)
}

func TestListTransfers(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	randomNum := util.RandomInt(0, 10)
	randomTransfers := make([]Transfer, randomNum)
	for i := 0; i < int(randomNum); i++ {
		randomTransfers[i] = createRandomTransfer(t, account1.ID, account2.ID)
	}

	params := ListTransfersParams{
		FromAccountID: sql.NullInt64{Int64: account1.ID, Valid: true},
		ToAccountID:   sql.NullInt64{Int64: account2.ID, Valid: true},
		Limit:         int32(randomNum),
	}

	transferList, err := testQueries.ListTransfers(context.Background(), params)

	require.NoError(t, err)
	require.NotEmpty(t, transferList)
	require.Len(t, transferList, int(randomNum))

	for i := 0; i < int(randomNum); i++ {
		require.Equal(t, randomTransfers[i].ID, transferList[i].ID)
		require.Equal(t, randomTransfers[i].ToAccountID, transferList[i].ToAccountID)
		require.Equal(t, randomTransfers[i].Amount, transferList[i].Amount)
	}
}
