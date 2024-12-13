package v1

import (
	"context"

	usecases "github.com/ppicom/newtonian/internal/application/use_cases"
)

// BankingServer implements the BankingServiceServer interface
type BankingServer struct {
	transferMoneyUseCase *usecases.TransferMoneyUseCase
}

// NewBankingServer creates a new BankingServer instance
func NewBankingServer(transferMoneyUseCase *usecases.TransferMoneyUseCase) *BankingServer {
	return &BankingServer{
		transferMoneyUseCase: transferMoneyUseCase,
	}
}

// TransferMoney handles money transfers between accounts
func (s *BankingServer) TransferMoney(ctx context.Context, req *TransferMoneyRequest) (*TransferMoneyResponse, error) {
	err := s.transferMoneyUseCase.Execute(
		req.GetFromAccountId(),
		req.GetToAccountId(),
		int(req.GetAmount()),
	)
	if err != nil {
		return nil, err
	}

	// For now, return a simple success response
	return &TransferMoneyResponse{
		Success: true,
		Message: "Transfer completed successfully",
	}, nil
}
