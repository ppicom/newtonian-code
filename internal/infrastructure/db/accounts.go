package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/ppicom/newtonian/internal/domain/banking"
	"github.com/redis/go-redis/v9"
)

const findAccountQuery = `SELECT id, balance FROM accounts WHERE id = ?`
const saveAccountQuery = `INSERT INTO accounts (id, balance) VALUES (?, ?) 
								ON DUPLICATE KEY UPDATE balance = ?`

type AccountRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func (r *AccountRepository) Find(id string) (*banking.Account, error) {
	if account, err := r.findInCache(id); err == nil {
		return account, nil
	}

	account, err := r.findInDatabase(id)
	if err != nil {
		return nil, err
	}

	r.updateCache(account)
	return account, nil
}

func (r *AccountRepository) findInCache(id string) (*banking.Account, error) {
	ctx := context.Background()
	data, err := r.redis.Get(ctx, "account:"+id).Bytes()
	if err != nil {
		return nil, err
	}

	var account banking.Account
	if err := json.Unmarshal(data, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) findInDatabase(id string) (*banking.Account, error) {
	var account banking.Account
	err := r.db.QueryRow(findAccountQuery, id).
		Scan(&account.ID, &account.Balance)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) Save(account *banking.Account) error {
	if err := r.saveToDatabase(account); err != nil {
		return err
	}

	r.updateCache(account)
	return nil
}

func (r *AccountRepository) saveToDatabase(account *banking.Account) error {
	_, err := r.db.Exec(saveAccountQuery, account.ID, account.Balance, account.Balance)
	return err
}

func (r *AccountRepository) updateCache(account *banking.Account) {
	ctx := context.Background()
	if accountJson, err := json.Marshal(account); err == nil {
		r.redis.Set(ctx, "account:"+account.ID, accountJson, 1*time.Hour)
	}
}

func NewAccountRepository(db *sql.DB, redis *redis.Client) *AccountRepository {
	return &AccountRepository{
		db:    db,
		redis: redis,
	}
}
