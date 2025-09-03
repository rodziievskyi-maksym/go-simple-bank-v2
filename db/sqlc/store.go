package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (s *Store) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	/*		//TxOptions holds the transaction options to be used in [DB.BeginTx].
	type TxOptions struct {
			//Isolation is the transaction isolation level.
			//If zero, the driver or database's default level is used.
		Isolation IsolationLevel
		ReadOnly  bool
	}*/
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		return err
	}

	queries := New(tx)
	if err = fn(queries); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %w, rollback err: %w", err, rbErr)
		}

		return err
	}

	return tx.Commit()
}

// TransferTxParams is the parameters that required to TransferTx transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"fromAccountId"`
	ToAccountID   int64 `json:"toAccountId"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of TransferTx transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromEntry   Entry    `json:"fromEntry"`
	ToEntry     Entry    `json:"toEntry"`
	FromAccount Account  `json:"fromAccount"`
	ToAccount   Account  `json:"toAccount"`
}

// TransferTx performs money from one account to another.
// It creates Transfer record, Entry records from and to, updates accounts balance within a single database transaction
func (s *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	if err := s.execTx(ctx, func(queries *Queries) error {
		var err error

		result.Transfer, err = queries.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// The best way to avoid database transaction deadlocks is to have consistent order of queries
		// Always update the smallest account id first.

		initPayload := NewPayload(ctx, queries, arg)

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = updateAccountBalance(initPayload)
			if err != nil {
				return err
			}
		} else {
			initPayload.swap()
			result.ToAccount, result.FromAccount, err = updateAccountBalance(initPayload)
			if err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return result, err
	}

	return result, nil
}

type UpdateAccountBalancePayload struct {
	ctx               context.Context
	queries           *Queries
	smallestAccountID int64
	amountToSmallest  int64
	biggestAccountID  int64
	amountToBiggest   int64
}

func NewPayload(ctx context.Context, q *Queries, arg TransferTxParams) *UpdateAccountBalancePayload {
	return &UpdateAccountBalancePayload{
		ctx:               ctx,
		queries:           q,
		smallestAccountID: arg.FromAccountID,
		amountToSmallest:  -arg.Amount,
		biggestAccountID:  arg.ToAccountID,
		amountToBiggest:   arg.Amount,
	}
}

func (p *UpdateAccountBalancePayload) swap() {
	smallestAccountID := p.smallestAccountID
	amountToSmallest := p.amountToSmallest

	p.smallestAccountID = p.biggestAccountID
	p.amountToSmallest = p.amountToBiggest
	p.biggestAccountID = smallestAccountID
	p.amountToBiggest = amountToSmallest
}

func updateAccountBalance(payload *UpdateAccountBalancePayload) (updatedAccount1, updatedAccount2 Account, err error) {
	if payload != nil {
		updatedAccount1, err = payload.queries.AddAccountBalance(payload.ctx, AddAccountBalanceParams{
			ID:     payload.smallestAccountID,
			Amount: payload.amountToSmallest,
		})
		if err != nil {
			return
		}

		updatedAccount2, err = payload.queries.AddAccountBalance(payload.ctx, AddAccountBalanceParams{
			ID:     payload.biggestAccountID,
			Amount: payload.amountToBiggest,
		})
		if err != nil {
			return
		}
	}

	return
}
