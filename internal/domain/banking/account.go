package banking

import "errors"

type Account struct {
	ID      string
	Balance int
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
	if err := Withdraw(from, amount); err != nil {
		return err
	}
	return Deposit(to, amount)
}
