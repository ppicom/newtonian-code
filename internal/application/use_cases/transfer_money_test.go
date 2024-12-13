package usecases_test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ppicom/newtonian/internal/infrastructure/db"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"

	usecases "github.com/ppicom/newtonian/internal/application/use_cases"
)

var (
	testDB    *sql.DB
	testRedis *redis.Client
)

func TestMain(m *testing.M) {
	// Setup test infrastructure
	var err error
	testDB, err = sql.Open("mysql", "test:testpass@tcp(localhost:3306)/banking_test")
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	// Wait for MySQL to be ready
	for i := 0; i < 30; i++ {
		err = testDB.Ping()
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		log.Fatalf("MySQL not ready after 30 seconds: %v", err)
	}

	// Create test table
	_, err = testDB.Exec(`
		CREATE TABLE IF NOT EXISTS accounts (
			id VARCHAR(255) PRIMARY KEY,
			balance INT NOT NULL
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create test table: %v", err)
	}

	testRedis = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Wait for Redis to be ready
	ctx := context.Background()
	for i := 0; i < 30; i++ {
		_, err = testRedis.Ping(ctx).Result()
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		log.Fatalf("Redis not ready after 30 seconds: %v", err)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	_, _ = testDB.Exec("DROP TABLE accounts")
	_ = testDB.Close()
	_ = testRedis.Close()

	os.Exit(code)
}

func setupTest(t *testing.T) (*usecases.TransferMoneyUseCase, func()) {
	t.Helper()

	// Clear data before each test
	_, err := testDB.Exec("DELETE FROM accounts")
	require.NoError(t, err)

	err = testRedis.FlushAll(context.Background()).Err()
	require.NoError(t, err)

	repo := db.NewAccountRepository(testDB, testRedis)
	useCase := usecases.NewTransferMoneyUseCase(repo)

	cleanup := func() {
		_, _ = testDB.Exec("DELETE FROM accounts")
		_ = testRedis.FlushAll(context.Background()).Err()
	}

	return useCase, cleanup
}

func createAccount(t *testing.T, id string, balance int) {
	t.Helper()
	_, err := testDB.Exec("INSERT INTO accounts (id, balance) VALUES (?, ?)", id, balance)
	require.NoError(t, err)
}

func getAccountBalance(t *testing.T, id string) int {
	t.Helper()
	var balance int
	err := testDB.QueryRow("SELECT balance FROM accounts WHERE id = ?", id).Scan(&balance)
	require.NoError(t, err)
	return balance
}

func TestTransferMoneyUseCase_Execute(t *testing.T) {
	tests := []struct {
		name            string
		fromID          string
		fromBalance     int
		toID            string
		toBalance       int
		amount          int
		expectedError   string
		expectedFromBal int
		expectedToBal   int
	}{
		{
			name:            "successful transfer",
			fromID:          "acc1",
			fromBalance:     100,
			toID:            "acc2",
			toBalance:       50,
			amount:          30,
			expectedFromBal: 70,
			expectedToBal:   80,
		},
		{
			name:          "insufficient balance",
			fromID:        "acc3",
			fromBalance:   100,
			toID:          "acc4",
			toBalance:     50,
			amount:        150,
			expectedError: "insufficient balance",
		},
		{
			name:          "invalid amount",
			fromID:        "acc5",
			fromBalance:   100,
			toID:          "acc6",
			toBalance:     50,
			amount:        0,
			expectedError: "invalid amount",
		},
		{
			name:          "account not found",
			fromID:        "non-existent",
			fromBalance:   0,
			toID:          "acc7",
			toBalance:     50,
			amount:        30,
			expectedError: "sql: no rows in result set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase, cleanup := setupTest(t)
			defer cleanup()

			// Setup test accounts
			if tt.fromID != "non-existent" {
				createAccount(t, tt.fromID, tt.fromBalance)
			}
			createAccount(t, tt.toID, tt.toBalance)

			// Execute transfer
			err := useCase.Execute(tt.fromID, tt.toID, tt.amount)

			// Verify results
			if tt.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)
			fromBalance := getAccountBalance(t, tt.fromID)
			toBalance := getAccountBalance(t, tt.toID)
			require.Equal(t, tt.expectedFromBal, fromBalance)
			require.Equal(t, tt.expectedToBal, toBalance)
		})
	}
}

func TestTransferMoneyUseCase_ConcurrentTransfers(t *testing.T) {
	useCase, cleanup := setupTest(t)
	defer cleanup()

	// Setup test accounts
	createAccount(t, "acc1", 1000)
	createAccount(t, "acc2", 1000)

	// Number of concurrent transfers
	numTransfers := 10
	transferAmount := 100

	// Create a channel to collect errors
	errChan := make(chan error, numTransfers*2)

	// Start concurrent transfers in both directions
	for i := 0; i < numTransfers; i++ {
		go func() {
			errChan <- useCase.Execute("acc1", "acc2", transferAmount)
		}()
		go func() {
			errChan <- useCase.Execute("acc2", "acc1", transferAmount)
		}()
	}

	// Collect all errors
	var errors []error
	for i := 0; i < numTransfers*2; i++ {
		if err := <-errChan; err != nil {
			errors = append(errors, err)
		}
	}

	// Verify results
	require.Empty(t, errors, "Expected no errors in concurrent transfers")

	// Final balances should be the same as initial ones
	acc1Balance := getAccountBalance(t, "acc1")
	acc2Balance := getAccountBalance(t, "acc2")
	require.Equal(t, 1000, acc1Balance, "Account 1 balance should remain unchanged after bidirectional transfers")
	require.Equal(t, 1000, acc2Balance, "Account 2 balance should remain unchanged after bidirectional transfers")
}
