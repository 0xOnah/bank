package service_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/onahvictor/bank/internal/config"
	mockdb "github.com/onahvictor/bank/internal/db/mock"
	"github.com/onahvictor/bank/internal/entity"
	"github.com/onahvictor/bank/internal/sdk/auth"
	"github.com/onahvictor/bank/internal/service"
	httptransport "github.com/onahvictor/bank/internal/transport/http"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type createUserTest struct {
	arg      entity.User
	password string
}

// arg user called at runtime
func (cr createUserTest) Matches(x any) bool {
	arg, ok := x.(entity.User)
	if !ok {
		return false
	}
	ok = auth.ComparePassword([]byte(arg.HashedPassword), cr.password)
	if !ok {
		return false
	}

	cr.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(cr.arg, arg)
}

func (cr createUserTest) String() string {
	return fmt.Sprintf("matches arg %v and password %v", cr.arg, cr.password)
}

func EqCreateUser(arg entity.User, password string) gomock.Matcher {
	return createUserTest{arg: arg, password: password}
}

func TestCreateUser(t *testing.T) {
	user, err := entity.NewUser("victor", "secret", "victor onah", "onahvictor@gmail.com")
	require.NoError(t, err)

	type TestCases []struct {
		name          string
		body          any
		buildStubs    func(repo *mockdb.MockUserRepository)
		checkResponse func(r *httptest.ResponseRecorder)
	}
	var testCases = TestCases{
		{
			name: "OK: Account Created",
			body: map[string]string{
				"username": user.Username,
				"password": "secret",
				"fullname": user.FullName,
				"email":    user.Email.String(),
				"currency": randomAccount().Currency,
			},
			buildStubs: func(repo *mockdb.MockUserRepository) {
				repo.EXPECT().CreateUser(gomock.Any(), EqCreateUser(user, "secret")).Times(1).Return(&user, nil)
			},
			checkResponse: func(r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, r.Code)
			},
		},
		{
			name: "Empty User: Not created",
			body: map[string]string{},
			buildStubs: func(repo *mockdb.MockUserRepository) {
				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, r.Code)
			},
		}, {
			name: "Internal server Error",
			body: map[string]string{
				"username": user.Username,
				"password": "secret",
				"fullname": user.FullName,
				"email":    user.Email.String(),
			},
			buildStubs: func(repo *mockdb.MockUserRepository) {
				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(nil, sql.ErrConnDone)
			},
			checkResponse: func(r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, r.Code)
			},
		}, {
			name: "Wrong Email",
			body: map[string]string{
				"username": user.Username,
				"password": "secret",
				"fullname": user.FullName,
				"email":    "invalid_email",
			},
			buildStubs: func(repo *mockdb.MockUserRepository) {
				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, r.Code)
			},
		},
	}

	for _, value := range testCases {
		t.Run(value.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			maker, err := auth.NewJWTMaker("123456789123456789123456789123456789123456789")
			require.NoError(t, err)

			accountRepo := mockdb.NewMockAccountRepository(ctrl)
			transferRepo := mockdb.NewMockTransferRepository(ctrl)
			userRepo := mockdb.NewMockUserRepository(ctrl)
			accountSvc := service.NewAccountService(accountRepo)
			accountHandler := httptransport.NewAccountHandler(accountSvc, maker)

			transferSvc := service.NewTransferService(transferRepo, accountRepo)
			transfHand := httptransport.NewTranserHandler(transferSvc, maker)

			usrSvc := service.NewUserService(userRepo, maker, config.Config{})
			userHand := httptransport.NewUserHandler(usrSvc, maker)

			router := httptransport.NewRouter(accountHandler, transfHand, userHand)

			data, err := json.Marshal(value.body)
			require.NoError(t, err)

			reader := bytes.NewReader([]byte(data))
			req, err := http.NewRequest(http.MethodPost, "/user", reader)
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)

			value.buildStubs(userRepo)
			rec := httptest.NewRecorder()
			router.Mux.ServeHTTP(rec, req)
			value.checkResponse(rec)
		})
	}
}
