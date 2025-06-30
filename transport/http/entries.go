package httptransport

import (
	"github.com/gin-gonic/gin"
	"github.com/onahvictor/bank/service"
)

type EntryHandler struct {
	entSvc *service.AccountService
}

func NewEntryHandler(svc *service.AccountService) *EntryHandler {
	return &EntryHandler{entSvc: svc}
}

func (a *EntryHandler) CreateEntry(ctx *gin.Context) {

}
