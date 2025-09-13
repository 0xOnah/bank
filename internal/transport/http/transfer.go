package httptransport

import (
	"context"
	"net/http"

	"github.com/0xOnah/bank/internal/entity"
	"github.com/0xOnah/bank/internal/sdk/auth"
	"github.com/0xOnah/bank/internal/sdk/util"
	"github.com/0xOnah/bank/internal/transport/sdk/errorutil"
	"github.com/0xOnah/bank/internal/transport/sdk/middleware"
	"github.com/gin-gonic/gin"
)

type TransferService interface {
	CreateTransferTX(ctx context.Context, arg entity.CreateTransferInput, username string, currency string) (*entity.TransferTxResult, error)
}
type TransferHandler struct {
	tranServ TransferService
	token    auth.Authenticator
}

func NewTranserHandler(svc TransferService, token auth.Authenticator) *TransferHandler {
	return &TransferHandler{tranServ: svc, token: token}
}

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gte=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (t *TransferHandler) MapAccountRoutes(r *gin.Engine) {
	r.POST("/transfer", middleware.Authenication(t.token), t.CreateTransfer)
}

func (t *TransferHandler) CreateTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.ClientIP()
	payload := ctx.MustGet(middleware.AuthorizationPayLoadKey).(*auth.Payload)

	arg := entity.CreateTransferInput{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	transfer, err := t.tranServ.CreateTransferTX(ctx.Request.Context(), arg, payload.Username, req.Currency)
	if err != nil {
		if appErr, ok := err.(*errorutil.AppError); ok {
			ctx.JSON(errorutil.MapErrorToHttpStatus(appErr), util.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfer)

}
