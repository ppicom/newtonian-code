package banking

type Account struct {
	Balance int
}

func Deposit(account *Account, amount int) error {
	account.Balance += amount
	return nil
}

func Withdraw(account *Account, amount int) error {
	account.Balance -= amount
	return nil
}

func Transfer(from *Account, to *Account, amount int) error {
	if err := Withdraw(from, amount); err != nil {
		return err
	}
	return Deposit(to, amount)
}
