package dbutils

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

/**
* Database Utilities - Helper Functions
**/

/**
* Accepts a function that expects a transaction argument and helps wrap the function call
* with the initiation, passing, and error checking of that transaction.
**/
func execTx(db *sqlx.DB, fn func(tx *sqlx.Tx) error) error {
	tx, err := db.Beginx()

	if err != nil {
		fmt.Println("Error when attempting to start transaction:", err)
		return fmt.Errorf("Error when attempting to start transaction: %s", err)
	}

	// call the function passed in and provide the transaction to it
	err = fn(tx)

	// function errored, rollback
	if err != nil {
		tx.Rollback()
	}

	// no error, safe to commit
	if err = tx.Commit(); err != nil {
		tx.Rollback() // rollback if commit fails
		return err
	}

	return nil
}
