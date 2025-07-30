package repo

import (
	"context"

	"github.com/0xOnah/bank/internal/db/sqlc"
	"github.com/0xOnah/bank/internal/entity"
)

type entryRepo struct {
	db sqlc.SQLStore
}

func toEntityEntry(e *sqlc.Entry) entity.Entry {
	return entity.Entry{
		ID:        e.ID,
		AccountID: e.AccountID,
		Amount:    e.Amount,
		CreatedAt: e.CreatedAt,
	}
}

func NewEntryRepo(db sqlc.SQLStore) *entryRepo {
	return &entryRepo{db: db}
}

func (r *entryRepo) CreateEntry(ctx context.Context, arg entity.CreateEntryInput) (entity.Entry, error) {
	result, err := r.db.CreateEntry(ctx, sqlc.CreateEntryParams{
		AccountID: arg.AccountID,
		Amount:    arg.Amount,
	})
	return toEntityEntry(result), err
}

func (r *entryRepo) GetEntry(ctx context.Context, id int64) (entity.Entry, error) {
	result, err := r.db.GetEntry(ctx, id)
	return toEntityEntry(result), err
}

func (r *entryRepo) ListEntries(ctx context.Context, arg entity.ListEntriesInput) ([]entity.Entry, error) {
	results, err := r.db.ListEntries(ctx, sqlc.ListEntriesParams{
		AccountID: arg.AccountID,
		Limit:     arg.Limit,
		Offset:    arg.Offset,
	})
	if err != nil {
		return nil, err
	}
	entries := make([]entity.Entry, 0, len(results))
	for _, e := range results {
		entries = append(entries, toEntityEntry(e))
	}
	return entries, nil
}
