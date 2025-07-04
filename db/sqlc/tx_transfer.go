package db

import (
	"context"
	"fmt"
)

// TransferTxParams input utama
type TransferTxParams struct {
	FromAccountID int64
	ToAccountID   int64
	Amount        int64
}

// TransferTxResult hasil output transaksi
type TransferTxResult struct {
	Transfer    Transfer
	FromAccount Account
	ToAccount   Account
}

// jalankan fungsi dalam transaksi SQL
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rollback err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTx menjalankan transfer atomik
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// 1. Buat transfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// 2. Update saldo dari pengirim
		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:      arg.FromAccountID,
			Balance: -arg.Amount,
		})
		if err != nil {
			return err
		}

		// 3. Update saldo ke penerima
		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:      arg.ToAccountID,
			Balance: arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
