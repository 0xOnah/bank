package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/0xOnah/bank/internal/sdk/auth"
	"github.com/0xOnah/bank/internal/sdk/util"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader     = "authorization"
	authorizationTypeBearer = "bearer"
	AuthorizationPayLoadKey = "authorization_payload"
)

func Authenication(payload auth.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authorizationHeader)
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse(errors.New("auth header not provided")))
			return
		}
		authParams := strings.Fields(authHeader)
		if len(authParams) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse(errors.New("malformed auth header")))
			return
		}

		authorizationType := strings.ToLower(authParams[0])
		if authorizationType != authorizationTypeBearer {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse(errors.New("malformed auth header")))
			return
		}

		accessToken := authParams[1]
		payload, err := payload.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse(err))
			return
		}

		ctx.Set(AuthorizationPayLoadKey, payload)
		ctx.Next()
	}

}
