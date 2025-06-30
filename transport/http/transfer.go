package httptransport

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/onahvictor/bank/entity"
	"github.com/onahvictor/bank/service"
	"github.com/onahvictor/bank/util"
)

type TransferService interface {
	CreateTransferTX(ctx context.Context, arg entity.CreateTransferInput, currency string) (*entity.TransferTxResult, error)
}
type TransferHandler struct {
	tranServ TransferService
}

func NewTranserHandler(svc TransferService) *TransferHandler {
	return &TransferHandler{tranServ: svc}
}

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gte=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (t *TransferHandler) MapAccountRoutes(r *gin.Engine) {
	r.POST("/transfer", t.CreateTransfer)
}

func (t *TransferHandler) CreateTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	arg := entity.CreateTransferInput{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	transfer, err := t.tranServ.CreateTransferTX(ctx.Request.Context(), arg, req.Currency)
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
