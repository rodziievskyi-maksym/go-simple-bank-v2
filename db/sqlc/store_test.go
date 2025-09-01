package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Printf(">> balances before transfer:\n Account 1 = %d \n Account 2 = %d \n", account1.Balance, account2.Balance)

	amount := int64(10)
	arg := TransferTxParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        amount,
	}

	concurrentTransactions := 2
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < concurrentTransactions; i++ {
		txName := fmt.Sprintf("transfer-%d", i)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, arg)

			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < concurrentTransactions; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.AccountID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.AccountID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		fmt.Printf(">> balances on [%d] TX:\n Account 1 = %d \n Account 2 = %d \n", i, fromAccount.Balance, toAccount.Balance)

		// check account balances diff
		//must be equal to amount variable -> account1.Balance = 761 - fromAccount.Balance = 751 (subtracted amount) = 10
		diffFromAccount := account1.Balance - fromAccount.Balance
		//must be equal to amount variable -> toAccount.Balance = 418 init value + amount (10 transferred from account 1) = 419 - 418 init value of 418 = 10
		diffToAccount := toAccount.Balance - account2.Balance
		require.Equal(t, diffFromAccount, diffToAccount)
		//check the validity of transfer account must still have positive balance to be able to process the transfer
		require.True(t, diffFromAccount > 0)

		// this interesting check shows us that we have right calculation on concurrent transactions
		//TODO: ask LLM what's going on here
		require.True(t, diffFromAccount%amount == 0) // amount, 2 * amount, 3 * amount

		//each concurrent transaction increate difference by the amount  1 - 10, 2 - 20, 3 - 30 and by dividing by amount
		//we've got the number from 1 to N (numbers on concurrent transactions)
		successfulIteration := int(diffFromAccount / amount)
		require.True(t, successfulIteration >= 1 && successfulIteration <= concurrentTransactions)
		require.NotContains(t, existed, successfulIteration)
		existed[successfulIteration] = true
	}

	//check final updated account balances
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount1)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount2)

	fmt.Printf(">> balances after transfer:\n Account 1 = %d \n Account 2 = %d \n", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance-amount*int64(concurrentTransactions), updatedAccount1.Balance)
	require.Equal(t, account2.Balance+amount*int64(concurrentTransactions), updatedAccount2.Balance)
}
