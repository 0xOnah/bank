package service

import (
	"context"
	"fmt"

	"github.com/0xOnah/bank/internal/entity"
	"github.com/0xOnah/bank/internal/transport/sdk/errorutil"
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
func (t *TransferService) validateAccount(ctx context.Context, accountId int64, currency string) (*entity.Account, error) {
	account, err := t.accountRepo.GetAccountByID(ctx, accountId)
	if err != nil {
		return nil, errorutil.NewAppError(errorutil.ErrNotFound, fmt.Sprintf("account Id=%d not found", accountId), err)
	}
	if account.Currency != currency {
		return nil, errorutil.NewAppError(errorutil.ErrBadRequest, fmt.Sprintf("account id=%d currency mismatch: %s vs %s", accountId, account.Currency, currency), nil)
	}
	return account, nil
}

// Todo: balance check
func (t *TransferService) CreateTransferTX(ctx context.Context, arg entity.CreateTransferInput, username string, currency string) (*entity.TransferTxResult, error) {
	//sameAccount
	if arg.FromAccountID == arg.ToAccountID {
		return nil, errorutil.NewAppError(errorutil.ErrInvalidInput, "cannot transfer to the same account", nil)
	}
	//from
	account, err := t.validateAccount(ctx, arg.FromAccountID, currency)
	if err != nil {
		return nil, err
	}
	if account.Owner != username {
		return nil, errorutil.NewAppError(errorutil.ErrUnauthorized, "you do not own this account", nil)
	}
	//to
	_, err = t.validateAccount(ctx, arg.ToAccountID, currency)
	if err != nil {
		return nil, err
	}
	//transfer
	tranfer, err := t.transferRepo.CreateTransferTX(ctx, arg)
	if err != nil {
		return nil, errorutil.NewAppError(errorutil.ErrInternal, "internal error", err)
	}

	return tranfer, nil
}
