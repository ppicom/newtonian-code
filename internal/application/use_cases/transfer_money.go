package usecases

import (
	"github.com/ppicom/newtonian/internal/domain/banking"
)

type AccountRepository interface {
	Find(id string) (*banking.Account, error)
	Save(account *banking.Account) error
}

type TransferMoneyUseCase struct {
	accountRepository AccountRepository
}

func (uc *TransferMoneyUseCase) Execute(from, to string, amount int) error {
	fromAccount, err := uc.accountRepository.Find(from)
	if err != nil {
		return err
	}

	toAccount, err := uc.accountRepository.Find(to)
	if err != nil {
		return err
	}

	return banking.Transfer(fromAccount, toAccount, amount)
}

func NewTransferMoneyUseCase(accountRepository AccountRepository) *TransferMoneyUseCase {
	return &TransferMoneyUseCase{accountRepository: accountRepository}
}
