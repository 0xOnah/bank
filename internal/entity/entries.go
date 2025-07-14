package entity

import "time"

type CreateEntryInput struct {
	AccountID int64
	Amount    int64
}

type ListEntriesInput struct {
	AccountID int64
	Limit     int32
	Offset    int32
}

type Entry struct {
	ID        int64     `json:"id"`
	AccountID int64     `json:"account_id"`
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
