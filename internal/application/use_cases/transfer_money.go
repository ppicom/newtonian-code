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

	if err := banking.Transfer(fromAccount, toAccount, amount); err != nil {
		return err
	}

	if err := uc.accountRepository.Save(fromAccount); err != nil {
		return err
	}

	return uc.accountRepository.Save(toAccount)
}

func NewTransferMoneyUseCase(accountRepository AccountRepository) *TransferMoneyUseCase {
	return &TransferMoneyUseCase{accountRepository: accountRepository}
}
