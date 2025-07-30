package httptransport

import (
	"time"

	"github.com/0xOnah/bank/internal/entity"
)

type AccountResp struct {
	ID       int64  `json:"id"`
	Owner    string `json:"owner"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

type Transfer struct {
	ID            int64     `json:"id"`
	FromAccountID int64     `json:"from_account_id"`
	ToAccountID   int64     `json:"to_account_id"`
	Amount        int64     `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}

type Entry struct {
	ID        int64     `json:"id"`
	AccountID int64     `json:"account_id"`
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

func toAccountsResp(acc []*entity.Account) []*AccountResp {
	var accounts []*AccountResp
	for _, v := range acc {
		account := AccountResp{
			ID:       v.ID,
			Owner:    v.Owner,
			Balance:  v.Balance,
			Currency: v.Currency,
		}
		accounts = append(accounts, &account)
	}
	return accounts
}
