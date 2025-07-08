package httptransport

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/onahvictor/bank/internal/entity"
	"github.com/onahvictor/bank/internal/sdk/auth"
	"github.com/onahvictor/bank/internal/service"
	"github.com/onahvictor/bank/internal/transport/middleware"
	"github.com/onahvictor/bank/internal/util"
)

type TransferService interface {
	CreateTransferTX(ctx context.Context, arg entity.CreateTransferInput, username string, currency string) (*entity.TransferTxResult, error)
}
type TransferHandler struct {
	tranServ TransferService
	token    auth.Auntenticator
}

func NewTranserHandler(svc TransferService, token auth.Auntenticator) *TransferHandler {
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
	payload := ctx.MustGet(middleware.AuthorizationPayLoadKey).(*auth.Payload)

	arg := entity.CreateTransferInput{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	transfer, err := t.tranServ.CreateTransferTX(ctx.Request.Context(), arg, payload.Username, req.Currency)
	if err != nil {
		if appErr, ok := err.(*service.AppError); ok {
			ctx.JSON(mapErrorToStatus(appErr), util.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfer)

}
