package banking

import (
	"errors"
	"sync"
)

type Account struct {
	ID      string
	Balance int
	mu      sync.Mutex
}

func (a *Account) Lock() {
	a.mu.Lock()
}

func (a *Account) Unlock() {
	a.mu.Unlock()
}

func Deposit(account *Account, amount int) error {
	if amount <= 0 {
		return errors.New("invalid amount")
	}

	account.Balance += amount
	return nil
}

func Withdraw(account *Account, amount int) error {
	if account.Balance < amount {
		return errors.New("insufficient balance")
	}

	if amount <= 0 {
		return errors.New("invalid amount")
	}

	account.Balance -= amount
	return nil
}

func Transfer(from *Account, to *Account, amount int) error {
	// Lock accounts in a consistent order to prevent deadlocks
	firstAccount, secondAccount := from, to
	if from.ID > to.ID {
		firstAccount, secondAccount = to, from
	}

	firstAccount.Lock()
	secondAccount.Lock()
	defer firstAccount.Unlock()
	defer secondAccount.Unlock()

	return transfer(from, to, amount)
}

// private helper function to perform the actual transfer
func transfer(from *Account, to *Account, amount int) error {
	if err := Withdraw(from, amount); err != nil {
		return err
	}
	return Deposit(to, amount)
}
