package httptransport

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/0xOnah/bank/internal/entity"
	"github.com/0xOnah/bank/internal/sdk/auth"
	"github.com/0xOnah/bank/internal/sdk/netutil"
	"github.com/0xOnah/bank/internal/service"
	"github.com/0xOnah/bank/internal/transport/sdk/errorutil"
	"github.com/0xOnah/bank/internal/sdk/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userService interface {
	CreateUser(ctx context.Context, cr service.CreateUserInput) (entity.User, error)
	Login(ctx context.Context, lg service.Logininput) (*service.AuthResult, error)
	RenewAccessToken(ctx context.Context, refreshToken string) (service.RenewAccessToken, error)
}

type UserHandler struct {
	UsrSvc userService
	auth   auth.Authenticator
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
func NewUserHandler(us userService, auth auth.Authenticator) *UserHandler {
	return &UserHandler{UsrSvc: us, auth: auth}
}
func (t *UserHandler) MapAccountRoutes(r *gin.Engine) {
	r.POST("/user", t.CreateUser)
	r.POST("/login", t.LoginAccount)
	r.POST("/token/renew_access", t.RenewAccessToken)
}

func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	var createUser userReq
	err := ctx.ShouldBindJSON(&createUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	usArg := service.CreateUserInput{
		Username: createUser.Username,
		Password: createUser.Password,
		Fullname: createUser.FullName,
		Email:    createUser.Email,
	}
	userValue, err := uh.UsrSvc.CreateUser(ctx.Request.Context(), usArg)
	if err != nil {
		if appErr, ok := err.(*errorutil.AppError); ok {
			ctx.JSON(errorutil.MapErrorToHttpStatus(appErr), util.ErrorResponse(err))
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
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	SessionID             uuid.UUID `json:"session_id"`
	User                  userResp  `json:"user"`
}

func (uh *UserHandler) LoginAccount(ctx *gin.Context) {
	var loginUser userLoginReq
	err := ctx.ShouldBindJSON(&loginUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	login := service.Logininput{
		Username:  loginUser.Username,
		Password:  loginUser.Password,
		ClientIP:  netutil.GetClientIP(ctx.Request),
		UserAgent: ctx.Request.UserAgent(),
	}
	data, err := uh.UsrSvc.Login(ctx.Request.Context(), login)
	if err != nil {
		if appErr, ok := err.(*errorutil.AppError); ok {
			ctx.JSON(errorutil.MapErrorToHttpStatus(appErr), util.ErrorResponse(err))
			slog.Info("Handled client error in CreateAccount",
				slog.Int("statusCode", int(appErr.Code)),
				slog.String("message", appErr.Message),
				slog.String("error", appErr.Err.Error()))
			return
		}
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))

		return
	}

	result := userLoginResp{
		AccessToken:           data.AccessToken,
		AccessTokenExpiresAt:  data.AccessTokenExpiresAt,
		RefreshToken:          data.RefreshToken,
		RefreshTokenExpiresAt: data.RefreshTokenExpiresAt,
		SessionID:             data.SessionID,
		User: userResp{
			Username: data.User.Username,
			Email:    data.User.Email.String(),
			FullName: data.User.FullName,
		},
	}
	ctx.JSON(http.StatusOK, result)
}

type renew struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RenewAccessTokenResp struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (uh *UserHandler) RenewAccessToken(ctx *gin.Context) {
	var refreshToken renew
	err := ctx.ShouldBindJSON(&refreshToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.ClientIP()
	accessToken, err := uh.UsrSvc.RenewAccessToken(ctx.Request.Context(), refreshToken.RefreshToken)
	if err != nil {
		if appErr, ok := err.(*errorutil.AppError); ok {
			ctx.JSON(errorutil.MapErrorToHttpStatus(appErr), util.ErrorResponse(err))
			slog.Info("Handled client error in CreateAccount",
				slog.Int("statusCode", int(appErr.Code)),
				slog.String("message", appErr.Message),
				slog.String("error", appErr.Err.Error()))
			return
		}
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))

		return
	}

	access := RenewAccessTokenResp{
		AccessToken:          accessToken.AccessToken,
		AccessTokenExpiresAt: accessToken.AccessTokenExpiresAt,
	}
	ctx.JSON(http.StatusOK, access)
}
