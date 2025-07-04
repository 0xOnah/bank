package httptransport

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/onahvictor/bank/internal/entity"
	"github.com/onahvictor/bank/internal/service"
	"github.com/onahvictor/bank/internal/util"
)

type AccountService interface {
	CreateAccount(ctx context.Context, input entity.CreateAccountInput) (*entity.Account, error)
	GetAccountByID(ctx context.Context, id int64) (*entity.Account, error)
	ListAccount(ctx context.Context, arg entity.ListAccountInput) ([]*entity.Account, error)
}
type AccountHandler struct {
	accSvc AccountService
}

func NewAccountHandler(svc *service.AccountService) *AccountHandler {
	return &AccountHandler{accSvc: svc}
}

func (a *AccountHandler) MapAccountRoutes(r *gin.Engine) {
	r.POST("/accounts", a.CreateAccount)
	r.GET("/accounts/:id", a.GetAccountByID)
	r.GET("/accounts", a.listAccount)

}

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

// create account
func (a *AccountHandler) CreateAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	account, err := a.accSvc.CreateAccount(ctx, entity.CreateAccountInput{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	})
	if err != nil {
		fmt.Println(err)
		var appErr *service.AppError
		if ok := errors.As(err, &appErr); ok {
			ctx.JSON(mapErrorToStatus(appErr), util.ErrorResponse(err))
			slog.Info("Handled client error in CreateAccount",
				slog.Int("statusCode", int(appErr.Code)),
				slog.String("message", appErr.Message),
				slog.String("error", appErr.Err.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))

		return
	}
	ctx.JSON(http.StatusOK, AccountResp{
		ID:       account.ID,
		Balance:  account.Balance,
		Owner:    account.Currency,
		Currency: account.Currency,
	})

}

type getAccountByID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (a *AccountHandler) GetAccountByID(ctx *gin.Context) {
	var req getAccountByID
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	account, err := a.accSvc.GetAccountByID(ctx.Request.Context(), req.ID)
	if err != nil {
		var appErr *service.AppError
		if ok := errors.As(err, &appErr); ok {
			statusCode := mapErrorToStatus(appErr)
			slog.Debug("Unexpected service error in GetAccountByID:",
				slog.Int("statusCode", statusCode),
				slog.String("message", appErr.Message),
				slog.Any("error", appErr.Err)) 
			ctx.JSON(statusCode, util.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, AccountResp{
		ID:       account.ID,
		Balance:  account.Balance,
		Owner:    account.Owner,
		Currency: account.Currency,
	})
}

type listAccountRequest struct {
	PageID   int64 `form:"page_id" binding:"required,min=1"`
	PageSize int64 `form:"page_size" binding:"required,min=5,max=10"`
}

func (a *AccountHandler) listAccount(ctx *gin.Context) {
	var arg listAccountRequest
	if err := ctx.ShouldBindQuery(&arg); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	accounts, err := a.accSvc.ListAccount(ctx.Request.Context(), entity.ListAccountInput{
		Limit:  int32(arg.PageSize),
		Offset: int32(arg.PageID-1) * int32(arg.PageSize),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	if len(accounts) == 0 {
		ctx.JSON(http.StatusOK, []entity.Account{})
		return
	}

	ctx.JSON(http.StatusOK, toAccountsResp(accounts))

}
