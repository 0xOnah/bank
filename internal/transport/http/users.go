package httptransport

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/onahvictor/bank/internal/entity"
	"github.com/onahvictor/bank/internal/sdk/auth"
	"github.com/onahvictor/bank/internal/service"
	"github.com/onahvictor/bank/internal/util"
)

type userService interface {
	CreateUser(ctx context.Context, username, email, password, fullname string) (entity.User, error)
	Login(ctx context.Context, username, password string) (*service.AuthResult, error)
}

type UserHandler struct {
	UsrSvc userService
	auth   auth.Auntenticator
}

type userReq struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"fullname" binding:"required"`
}

type userResp struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	FullName string `json:"fullname" binding:"required"`
}

func (ur userReq) ToUserResponse(u entity.User) userResp {
	return userResp{
		Username: u.Username,
		Email:    u.Email.String(),
		FullName: u.FullName,
	}
}
func NewUserHandler(us userService, auth auth.Auntenticator) *UserHandler {
	return &UserHandler{UsrSvc: us, auth: auth}
}
func (t *UserHandler) MapAccountRoutes(r *gin.Engine) {
	r.POST("/user", t.CreateUser)
	r.POST("/login", t.LoginAccount)
}

func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	var createUser userReq
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

type userLoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type userLoginResp struct {
	AccessToken string   `json:"access_token"`
	User        userResp `json:"user"`
}

func (uh *UserHandler) LoginAccount(ctx *gin.Context) {
	var loginUser userLoginReq
	err := ctx.ShouldBindJSON(&loginUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	data, err := uh.UsrSvc.Login(ctx.Request.Context(), loginUser.Username, loginUser.Password)
	if err != nil {
		if appErr, ok := err.(*service.AppError); ok {
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

	result := userLoginResp{
		AccessToken: data.Token,
		User: userResp{
			Username: data.User.Username,
			Email:    data.User.Email.String(),
			FullName: data.User.FullName,
		},
	}
	ctx.JSON(http.StatusOK, result)
}
