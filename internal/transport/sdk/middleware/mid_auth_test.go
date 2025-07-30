package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/0xOnah/bank/internal/sdk/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestAuth(t *testing.T) {
	testAuth := []struct {
		name          string
		setupAuth     func(r *http.Request, tokenMaker auth.Authenticator)
		checkResponse func(record *httptest.ResponseRecorder)
	}{
		{
			name: "OK Valid token",
			setupAuth: func(r *http.Request, tokenMaker auth.Authenticator) {
				token, _, err := tokenMaker.GenerateToken("user", time.Minute*15)
				require.NoError(t, err)
				r.Header.Set("Authorization", fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			checkResponse: func(record *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, record.Code)
			},
		}, {
			name: "No Authorization",
			setupAuth: func(r *http.Request, tokenMaker auth.Authenticator) {

			},
			checkResponse: func(record *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, record.Code)
			},
		}, {
			name: "Expired token",
			setupAuth: func(r *http.Request, tokenMaker auth.Authenticator) {
				token, _, err := tokenMaker.GenerateToken("user", -time.Minute*15)
				require.NoError(t, err)
				r.Header.Set("Authorization", fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			checkResponse: func(record *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, record.Code)
			},
		}, {
			name: "Invalid AuthorizationFormat",
			setupAuth: func(r *http.Request, tokenMaker auth.Authenticator) {
				token, _, err := tokenMaker.GenerateToken("user", time.Minute*15)
				require.NoError(t, err)
				r.Header.Set("", fmt.Sprintf("%s %s", authorizationTypeBearer, token))
			},
			checkResponse: func(record *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, record.Code)
			},
		},
	}
	for i := range testAuth {
		tc := testAuth[i]
		t.Run(tc.name, func(t *testing.T) {
			maker, err := auth.NewJWTMaker("helloworldhelloworldhelloworldhelloworld")
			require.NoError(t, err)

			router := gin.New()
			router.GET("/auth", Authenication(maker), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{})
			})

			rec := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/auth", nil)
			require.NoError(t, err)

			tc.setupAuth(req, maker)
			router.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}
