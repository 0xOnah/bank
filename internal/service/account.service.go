package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/onahvictor/bank/internal/db/repo"
	"github.com/onahvictor/bank/internal/entity"
)

type AccountRepository interface {
	AddAccountBalance(ctx context.Context, arg entity.AddAccountBalanceInput) (*entity.Account, error)
	CreateAccount(ctx context.Context, arg entity.CreateAccountInput) (*entity.Account, error)
	DeleteAccount(ctx context.Context, id int64) error
	GetAccountByID(ctx context.Context, id int64) (*entity.Account, error)
	GetAccountForUpdate(ctx context.Context, id int64) (*entity.Account, error)
	ListAccount(ctx context.Context, arg entity.ListAccountInput) ([]*entity.Account, error)
	UpdateAccount(ctx context.Context, arg entity.UpdateAccountInput) (*entity.Account, error)
}

type AccountService struct {
	accountRepo AccountRepository
}

func NewAccountService(repo AccountRepository) *AccountService {
	return &AccountService{accountRepo: repo}
}

func (a *AccountService) CreateAccount(ctx context.Context, arg entity.CreateAccountInput) (*entity.Account, error) {
	account, err := a.accountRepo.CreateAccount(ctx, arg)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrDuplicateAccountCurrency):
			return nil, NewAppError(ErrBadRequest, fmt.Sprintf("an account with currency %s already exists for this user", arg.Currency), err)
		case errors.Is(err, repo.ErrUserNotExist):
			return nil, NewAppError(ErrBadRequest, "this user does not exist", err)
		default:
			return nil, NewAppError(ErrInternal, "internal server error", err)
		}
	}
	return account, nil
}

func (a *AccountService) GetAccountByID(ctx context.Context, id int64) (*entity.Account, error) {
	account, err := a.accountRepo.GetAccountByID(ctx, id)
	if err != nil {
		if errors.Is(err, repo.ErrRecordNotFound) {
			return nil, NewAppError(ErrNotFound, fmt.Sprintf("account %d not found", id), err)
		}
		return nil, NewAppError(ErrUnknown, "failed to retrieve account due to an unexpected issue", err) // `err` here is the original repo error
	}
	return account, nil
}

func (a *AccountService) ListAccount(ctx context.Context, arg entity.ListAccountInput) ([]*entity.Account, error) {
	accounts, err := a.accountRepo.ListAccount(ctx, arg)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
