package service

import (
	"context"

	"github.com/0xOnah/bank/internal/entity"
	"github.com/0xOnah/bank/internal/transport/sdk/errorutil"
)

type EntryRepository interface {
	CreateEntry(ctx context.Context, arg entity.CreateEntryInput) (entity.Entry, error)
	GetEntry(ctx context.Context, id int64) (entity.Entry, error)
	ListEntries(ctx context.Context, arg entity.ListEntriesInput) ([]entity.Entry, error)
}

type EntryService struct {
	entryRepo EntryRepository
}

func NewEntryService(repo EntryRepository) *EntryService {
	return &EntryService{entryRepo: repo}
}
func (e *EntryService) CreateEntry(ctx context.Context, arg entity.CreateEntryInput) (entity.Entry, error) {
	_, err := e.entryRepo.CreateEntry(ctx, arg)
	if err != nil {
		return entity.Entry{}, errorutil.NewAppError(errorutil.ErrInternal, "internal error", err)
	}
	return entity.Entry{}, nil
}
