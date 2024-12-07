package db

import (
	"database/sql"

	"github.com/ppicom/newtonian/internal/domain/banking"
)

const findAccountQuery = `SELECT id, balance FROM accounts WHERE id = ?`
const saveAccountQuery = `INSERT INTO accounts (id, balance) VALUES (?, ?) 
								ON DUPLICATE KEY UPDATE balance = ?`

type SqlAccountRepository struct {
	db *sql.DB
}

func (r *SqlAccountRepository) Find(id string) (*banking.Account, error) {
	var account banking.Account
	err := r.db.QueryRow(findAccountQuery, id).
		Scan(&account.ID, &account.Balance)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *SqlAccountRepository) Save(account *banking.Account) error {
	_, err := r.db.Exec(saveAccountQuery, account.ID, account.Balance, account.Balance)
	return err
}

func NewSqlAccountRepository(db *sql.DB) *SqlAccountRepository {
	return &SqlAccountRepository{db: db}
}
