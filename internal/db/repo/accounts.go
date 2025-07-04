package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/onahvictor/bank/internal/db/sqlc"
	"github.com/onahvictor/bank/internal/entity"
)

type accountRepo struct {
	db *sqlc.SQLStore
}

func toEntityAccount(a *sqlc.Account) *entity.Account {
	return &entity.Account{
		ID:        a.ID,
		Owner:     a.Owner,
		Balance:   a.Balance,
		Currency:  a.Currency,
		CreatedAt: a.CreatedAt,
	}
}

func NewAccountRepo(db *sqlc.SQLStore) *accountRepo {
	return &accountRepo{db: db}
}

func (r *accountRepo) AddAccountBalance(ctx context.Context, arg entity.AddAccountBalanceInput) (*entity.Account, error) {
	result, err := r.db.AddAccountBalance(ctx, sqlc.AddAccountBalanceParams{
		Amount: arg.Amount,
		ID:     arg.ID,
	})
	return toEntityAccount(result), err
}

func (r *accountRepo) CreateAccount(ctx context.Context, arg entity.CreateAccountInput) (*entity.Account, error) {
	result, err := r.db.CreateAccount(ctx, sqlc.CreateAccountParams{
		Owner:    arg.Owner,
		Balance:  arg.Balance,
		Currency: arg.Currency,
	})

	if err != nil {
		if errValue, ok := err.(*pq.Error); ok {
			switch errValue.Code.Name() {
			case "unique_violation":
				return nil, ErrDuplicateAccountCurrency
			case "foreign_key_violation":
				return nil, ErrUserNotExist
			}
		}
		return nil, err
	}
	return toEntityAccount(result), nil
}

func (r *accountRepo) DeleteAccount(ctx context.Context, id int64) error {
	return r.db.DeleteAccount(ctx, id)
}

func (r *accountRepo) GetAccountByID(ctx context.Context, id int64) (*entity.Account, error) {
	result, err := r.db.GetAccount(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return toEntityAccount(result), err
}

func (r *accountRepo) GetAccountForUpdate(ctx context.Context, id int64) (*entity.Account, error) {
	result, err := r.db.GetAccountForUpdate(ctx, id)

	return toEntityAccount(result), err
}

func (r *accountRepo) ListAccount(ctx context.Context, arg entity.ListAccountInput) ([]*entity.Account, error) {
	results, err := r.db.ListAccount(ctx, sqlc.ListAccountParams{
		Limit:  arg.Limit,
		Offset: arg.Offset,
	})
	if err != nil {
		return nil, err
	}
	accounts := make([]*entity.Account, 0, len(results))
	for _, a := range results {
		accounts = append(accounts, toEntityAccount(a))
	}
	return accounts, nil
}

func (r *accountRepo) UpdateAccount(ctx context.Context, arg entity.UpdateAccountInput) (*entity.Account, error) {
	result, err := r.db.UpdateAccount(ctx, sqlc.UpdateAccountParams{
		ID:      arg.ID,
		Balance: arg.Balance,
	})
	return toEntityAccount(result), err
}
