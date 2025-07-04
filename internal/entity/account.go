package entity

import "time"

type Users struct{
	
}

type Account struct {
	ID        int64     `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"-"`
}

type AddAccountBalanceInput struct {
	ID     int64
	Amount int64
}

type CreateAccountInput struct {
	Owner    string
	Balance  int64
	Currency string
}

type UpdateAccountInput struct {
	ID      int64
	Balance int64
}

type ListAccountInput struct {
	Limit  int32
	Offset int32
}
