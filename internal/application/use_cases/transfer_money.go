package usecases

import (
	"database/sql"
	"sync"

	"github.com/ppicom/newtonian/internal/domain/banking"
)

type AccountRepository interface {
	Find(tx *sql.Tx, id string) (*banking.Account, error)
	Save(tx *sql.Tx, account *banking.Account) error
	BeginTx() (*sql.Tx, error)
	CommitTx(tx *sql.Tx) error
	RollbackTx(tx *sql.Tx) error
}

type TransferMoneyUseCase struct {
	accountRepository AccountRepository
	mu                sync.Mutex
}

func (uc *TransferMoneyUseCase) Execute(from, to string, amount int) error {
	// Lock the use case to guarantee concurrent transfers are serialized and happen in order
	uc.mu.Lock()
	defer uc.mu.Unlock()

	tx, err := uc.accountRepository.BeginTx()
	if err != nil {
		return err
	}

	fromAccount, err := uc.accountRepository.Find(tx, from)
	if err != nil {
		uc.accountRepository.RollbackTx(tx)
		return err
	}

	toAccount, err := uc.accountRepository.Find(tx, to)
	if err != nil {
		uc.accountRepository.RollbackTx(tx)
		return err
	}

	if err := banking.Transfer(fromAccount, toAccount, amount); err != nil {
		uc.accountRepository.RollbackTx(tx)
		return err
	}

	if err := uc.accountRepository.Save(tx, fromAccount); err != nil {
		uc.accountRepository.RollbackTx(tx)
		return err
	}

	if err := uc.accountRepository.Save(tx, toAccount); err != nil {
		uc.accountRepository.RollbackTx(tx)
		return err
	}

	return uc.accountRepository.CommitTx(tx)
}

func NewTransferMoneyUseCase(accountRepository AccountRepository) *TransferMoneyUseCase {
	return &TransferMoneyUseCase{accountRepository: accountRepository}
}
