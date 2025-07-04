package service

import (
	"context"
	"fmt"

	"github.com/onahvictor/bank/internal/entity"
)

type TransferRepository interface {
	CreateTransfer(ctx context.Context, arg entity.CreateTransferInput) (*entity.Transfer, error)
	GetTransfer(ctx context.Context, id int64) (*entity.Transfer, error)
	ListTransfers(ctx context.Context, arg entity.ListTransfersInput) ([]*entity.Transfer, error)
	CreateTransferTX(ctx context.Context, arg entity.CreateTransferInput) (*entity.TransferTxResult, error)
}

type TransferService struct {
	transferRepo TransferRepository
	accountRepo  AccountRepository
}

func NewTransferService(transRepo TransferRepository, accountRepo AccountRepository) *TransferService {
	return &TransferService{
		accountRepo:  accountRepo,
		transferRepo: transRepo,
	}
}
func (t *TransferService) validateAccount(ctx context.Context, accountId int64, currency string) error {
	account, err := t.accountRepo.GetAccountByID(ctx, accountId)
	if err != nil {
		return NewAppError(ErrNotFound, fmt.Sprintf("account Id=%d not found", accountId), err)
	}
	if account.Currency != currency {
		return NewAppError(ErrBadRequest, fmt.Sprintf("account id=%d currency mismatch: %s vs %s", accountId, account.Currency, currency), nil)
	}
	return nil
}

// Todo: balance check
func (t *TransferService) CreateTransferTX(ctx context.Context, arg entity.CreateTransferInput, currency string) (*entity.TransferTxResult, error) {
	//sameAccount
	if arg.FromAccountID == arg.ToAccountID {
		return nil, NewAppError(ErrInvalidInput, "cannot transfer to the same account", nil)
	}
	//from
	err := t.validateAccount(ctx, arg.FromAccountID, currency)
	if err != nil {
		return nil, err
	}
	//to
	err = t.validateAccount(ctx, arg.ToAccountID, currency)
	if err != nil {
		return nil, err
	}
	//transfer
	tranfer, err := t.transferRepo.CreateTransferTX(ctx, arg)
	if err != nil {
		return nil, NewAppError(ErrInternal, "internal error", err)
	}

	return tranfer, nil
}
