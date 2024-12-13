package banking_test

import (
	"testing"

	"github.com/ppicom/newtonian/internal/domain/banking"
)

func TestDeposit(t *testing.T) {
	tests := []struct {
		name          string
		initialAmount int
		deposit       int
		wantBalance   int
		wantError     bool
	}{
		{
			name:          "valid deposit",
			initialAmount: 100,
			deposit:       50,
			wantBalance:   150,
			wantError:     false,
		},
		{
			name:          "zero deposit",
			initialAmount: 100,
			deposit:       0,
			wantBalance:   100,
			wantError:     true,
		},
		{
			name:          "negative deposit",
			initialAmount: 100,
			deposit:       -50,
			wantBalance:   100,
			wantError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := &banking.Account{Balance: tt.initialAmount}
			err := banking.Deposit(account, tt.deposit)

			if (err != nil) != tt.wantError {
				t.Errorf("Deposit() error = %v, wantError %v", err, tt.wantError)
			}

			if account.Balance != tt.wantBalance {
				t.Errorf("Balance = %v, want %v", account.Balance, tt.wantBalance)
			}
		})
	}
}

func TestWithdraw(t *testing.T) {
	tests := []struct {
		name          string
		initialAmount int
		withdraw      int
		wantBalance   int
		wantError     bool
	}{
		{
			name:          "valid withdrawal",
			initialAmount: 100,
			withdraw:      50,
			wantBalance:   50,
			wantError:     false,
		},
		{
			name:          "zero withdrawal",
			initialAmount: 100,
			withdraw:      0,
			wantBalance:   100,
			wantError:     true,
		},
		{
			name:          "negative withdrawal",
			initialAmount: 100,
			withdraw:      -50,
			wantBalance:   100,
			wantError:     true,
		},
		{
			name:          "insufficient balance",
			initialAmount: 100,
			withdraw:      150,
			wantBalance:   100,
			wantError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account := &banking.Account{Balance: tt.initialAmount}
			err := banking.Withdraw(account, tt.withdraw)

			if (err != nil) != tt.wantError {
				t.Errorf("Withdraw() error = %v, wantError %v", err, tt.wantError)
			}

			if account.Balance != tt.wantBalance {
				t.Errorf("Balance = %v, want %v", account.Balance, tt.wantBalance)
			}
		})
	}
}

func TestTransfer(t *testing.T) {
	tests := []struct {
		name              string
		fromInitialAmount int
		toInitialAmount   int
		transferAmount    int
		wantFromBalance   int
		wantToBalance     int
		wantError         bool
	}{
		{
			name:              "valid transfer",
			fromInitialAmount: 100,
			toInitialAmount:   50,
			transferAmount:    30,
			wantFromBalance:   70,
			wantToBalance:     80,
			wantError:         false,
		},
		{
			name:              "insufficient balance",
			fromInitialAmount: 100,
			toInitialAmount:   50,
			transferAmount:    150,
			wantFromBalance:   100,
			wantToBalance:     50,
			wantError:         true,
		},
		{
			name:              "zero transfer",
			fromInitialAmount: 100,
			toInitialAmount:   50,
			transferAmount:    0,
			wantFromBalance:   100,
			wantToBalance:     50,
			wantError:         true,
		},
		{
			name:              "negative transfer",
			fromInitialAmount: 100,
			toInitialAmount:   50,
			transferAmount:    -30,
			wantFromBalance:   100,
			wantToBalance:     50,
			wantError:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			from := &banking.Account{Balance: tt.fromInitialAmount}
			to := &banking.Account{Balance: tt.toInitialAmount}

			err := banking.Transfer(from, to, tt.transferAmount)

			if (err != nil) != tt.wantError {
				t.Errorf("Transfer() error = %v, wantError %v", err, tt.wantError)
			}

			if from.Balance != tt.wantFromBalance {
				t.Errorf("From Balance = %v, want %v", from.Balance, tt.wantFromBalance)
			}

			if to.Balance != tt.wantToBalance {
				t.Errorf("To Balance = %v, want %v", to.Balance, tt.wantToBalance)
			}
		})
	}
}
