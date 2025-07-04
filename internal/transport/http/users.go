package httptransport

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/onahvictor/bank/internal/entity"
	"github.com/onahvictor/bank/internal/service"
	"github.com/onahvictor/bank/internal/util"
)

type userService interface {
	CreateUser(ctx context.Context, username, email, password, fullname string) (entity.User, error)
}

type UserHandler struct {
	UsrSvc userService
}

type user struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"fullname" binding:"required"`
}

func (ur user) ToUserResponse(u entity.User) user {
	return user{
		Username: u.Username,
		Email:    u.Email.String(),
		FullName: u.FullName,
	}
}
func NewUserHandler(us userService) *UserHandler {
	return &UserHandler{UsrSvc: us}
}
func (t *UserHandler) MapAccountRoutes(r *gin.Engine) {
	r.POST("/user", t.CreateUser)
}
func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	var createUser user
	err := ctx.ShouldBindJSON(&createUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	userValue, err := uh.UsrSvc.CreateUser(ctx.Request.Context(), createUser.Username, createUser.Email, createUser.Password, createUser.FullName)
	if err != nil {
		if appErr, ok := err.(*service.AppError); ok {
			ctx.JSON(mapErrorToStatus(appErr), util.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, createUser.ToUserResponse(userValue))

}
