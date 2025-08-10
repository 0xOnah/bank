package httptransport

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/0xOnah/bank/internal/entity"
	"github.com/0xOnah/bank/internal/sdk/auth"
	"github.com/0xOnah/bank/internal/service"
	"github.com/0xOnah/bank/internal/transport/sdk/errorutil"
	"github.com/0xOnah/bank/internal/transport/sdk/middleware"
	"github.com/0xOnah/bank/internal/sdk/util"
	"github.com/gin-gonic/gin"
)

type AccountService interface {
	CreateAccount(ctx context.Context, input entity.CreateAccountInput) (*entity.Account, error)
	GetAccountByID(ctx context.Context, username string, id int64) (*entity.Account, error)
	ListAccount(ctx context.Context, arg entity.ListAccountInput) ([]*entity.Account, error)
}
type AccountHandler struct {
	accSvc AccountService
	token  auth.Authenticator
}

func NewAccountHandler(svc *service.AccountService, token auth.Authenticator) *AccountHandler {
	return &AccountHandler{accSvc: svc, token: token}
}

func (a *AccountHandler) MapAccountRoutes(r *gin.Engine) {
	r.POST("/accounts", middleware.Authenication(a.token), a.CreateAccount)
	r.GET("/accounts/:id", middleware.Authenication(a.token), a.GetAccountByID)
	r.GET("/accounts", middleware.Authenication(a.token), a.listAccount)
}

type CreateAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

// create account
func (a *AccountHandler) CreateAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	payload := ctx.MustGet(middleware.AuthorizationPayLoadKey).(*auth.Payload)

	account, err := a.accSvc.CreateAccount(ctx, entity.CreateAccountInput{
		Owner:    payload.Username,
		Currency: req.Currency,
		Balance:  0,
	})
	if err != nil {
		var appErr *errorutil.AppError
		if ok := errors.As(err, &appErr); ok {
			ctx.JSON(errorutil.MapErrorToHttpStatus(appErr), util.ErrorResponse(err))
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
	payload := ctx.MustGet(middleware.AuthorizationPayLoadKey).(*auth.Payload)

	account, err := a.accSvc.GetAccountByID(ctx.Request.Context(), payload.Username, req.ID)
	if err != nil {
		var appErr *errorutil.AppError
		if ok := errors.As(err, &appErr); ok {
			statusCode := errorutil.MapErrorToHttpStatus(appErr)
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

	payload := ctx.MustGet(middleware.AuthorizationPayLoadKey).(*auth.Payload)
	accounts, err := a.accSvc.ListAccount(ctx.Request.Context(), entity.ListAccountInput{
		User:   payload.Username,
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
