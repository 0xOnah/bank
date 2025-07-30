package repo

import (
	"context"

	"github.com/0xOnah/bank/internal/db/sqlc"
	"github.com/0xOnah/bank/internal/entity"
)

type transferRepo struct {
	db *sqlc.SQLStore
}

func NewTransferRepo(db *sqlc.SQLStore) *transferRepo {
	return &transferRepo{db: db}
}

func NewAccountResp(acc *sqlc.Account) *entity.Account {
	if acc == nil {
		return nil
	}
	return &entity.Account{
		ID:       acc.ID,
		Owner:    acc.Owner,
		Balance:  acc.Balance,
		Currency: acc.Currency,
	}
}
func NewTransfResp(trans *sqlc.Transfer) *entity.Transfer {
	if trans == nil {
		return nil
	}
	return &entity.Transfer{
		ID:            trans.ID,
		FromAccountID: trans.FromAccountID,
		ToAccountID:   trans.ToAccountID,
		Amount:        trans.Amount,
		CreatedAt:     trans.CreatedAt,
	}
}

func NewEntryResp(entry *sqlc.Entry) *entity.Entry {
	if entry == nil {
		return nil
	}
	return &entity.Entry{
		ID:        entry.ID,
		AccountID: entry.AccountID,
		Amount:    entry.Amount,
		CreatedAt: entry.CreatedAt,
	}
}

func NewTransferTxResponse(txResult *sqlc.TransferTxResult) *entity.TransferTxResult {
	return &entity.TransferTxResult{
		Transfer:    NewTransfResp(txResult.Transfer),
		FromAccount: NewAccountResp(txResult.FromAccount),
		ToAccount:   NewAccountResp(txResult.ToAccount),
		FromEntry:   NewEntryResp(txResult.FromEntry),
		ToEntry:     NewEntryResp(txResult.ToEntry),
	}
}

func (r *transferRepo) CreateTransferTX(ctx context.Context, arg entity.CreateTransferInput) (*entity.TransferTxResult, error) {
	result, err := r.db.TransferTx(ctx, sqlc.TransferTxParams{
		FromAccountID: arg.FromAccountID,
		ToAccountID:   arg.ToAccountID,
		Amount:        arg.Amount,
	})
	return NewTransferTxResponse(result), err
}

func (r *transferRepo) CreateTransfer(ctx context.Context, arg entity.CreateTransferInput) (*entity.Transfer, error) {
	result, err := r.db.CreateTransfer(ctx, sqlc.CreateTransferParams{
		FromAccountID: arg.FromAccountID,
		ToAccountID:   arg.ToAccountID,
		Amount:        arg.Amount,
	})
	return toEntityTransfer(result), err
}

func (r *transferRepo) GetTransfer(ctx context.Context, id int64) (*entity.Transfer, error) {
	result, err := r.db.GetTransfer(ctx, id)
	return toEntityTransfer(result), err
}

func (r *transferRepo) ListTransfers(ctx context.Context, arg entity.ListTransfersInput) ([]*entity.Transfer, error) {
	results, err := r.db.ListTransfers(ctx, sqlc.ListTransfersParams{
		FromAccountID: arg.FromAccountID,
		ToAccountID:   arg.ToAccountID,
		Limit:         arg.Limit,
		Offset:        arg.Offset,
	})
	if err != nil {
		return nil, err
	}
	transfers := make([]*entity.Transfer, 0, len(results))
	for _, t := range results {
		transfers = append(transfers, toEntityTransfer(t))
	}
	return transfers, nil
}

func toEntityTransfer(t *sqlc.Transfer) *entity.Transfer {
	return &entity.Transfer{
		ID:            t.ID,
		FromAccountID: t.FromAccountID,
		ToAccountID:   t.ToAccountID,
		Amount:        t.Amount,
		CreatedAt:     t.CreatedAt,
	}
}
