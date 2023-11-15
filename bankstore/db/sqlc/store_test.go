package db

import (
	// "context"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println("==============================")
	fmt.Println("FromAccID: ", account1.ID, "ToAccID:", account2.ID)
	fmt.Println("Before: ", account1.Balance, "---", account2.Balance)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)
	// run n concurrent transfer transaction
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
		time.Sleep(1 * time.Second) 
	}

	// Check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		// check transfer
		result := <-results
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// TODO: check accounts' balance
	}
	accountFrom, _ := store.GetAccount(context.Background(), account1.ID)
	accountTo, _ := store.GetAccount(context.Background(), account2.ID)
	require.Equal(t, account1.Balance, accountFrom.Balance+amount*int64(n))
	require.Equal(t, account2.Balance, accountTo.Balance-amount*int64(n))

	fmt.Println("==============================")
	fmt.Println("FromAccID: ", accountFrom.ID, "ToAccID:", accountTo.ID)
	fmt.Println("After: ", accountFrom.Balance, "-->", accountTo.Balance, "Amount: ", amount*int64(n))
	fmt.Println("==============================")
}
