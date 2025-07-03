package db

import (
	"context"
	"fmt"
)

type TrasferTxParams struct {
	FromAccountID int64
	ToAccountID   int64
	Amount        int64
}

type TrasferTxResult struct {
	Transfer Transfer
}

func (store *SQLStore) TrasferTx(ctx context.Context,
	arg TrasferTxParams) (TrasferTxResult, error) {
	var result TrasferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		transfer, err := q.CreateTransfer(ctx,
			CreateTransferParams{
				FromAccountID: arg.FromAccountID,
				ToAccountID:   arg.ToAccountID,
				Amount:        arg.Amount,
			})
		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			_, err =
				addMoney(ctx, q, arg.FromAccountID, -arg.Amount)
			if err != nil {
				return err
			}
		}

		result.Transfer = transfer
		return nil
	})
	return result, err
}

func (store *SQLStore) execTx(ctx context.Context,
	fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rollback err: %v",
				err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

func addMoney(ctx context.Context, q *Queries,
	accountID int64, amount int64) (Account, error) {
	return q.AddAccountBalance(ctx,
		AddAccountBalanceParams{
			ID:      accountID,
			Balance: amount,
		})
}
