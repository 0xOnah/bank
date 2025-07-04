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
		checkResponse func(t *testing.T, r *httptest.ResponseRecorder)
	}
	var testCases = TestCases{
		{name: "OK: Account Created",
			body: map[string]string{
				"username": "victor",
				"password": "secret",
				"fullname": "victor onah",
				"email":    "onahvictor@gmail.com"},
			buildStubs: func(repo *mockdb.MockUserRepository) {
				repo.EXPECT().CreateUser(gomock.Any(), EqCreateUser(user, "secret")).Times(1).Return(&user, nil)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, r.Code)
			},
		},
		{
			name: "Empty User: Not created",
			body: map[string]string{},
			buildStubs: func(repo *mockdb.MockUserRepository) {
				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, r.Code)
			},
		}, {
			name: "Internal server Error",
			body: map[string]string{
				"username": "victor",
				"password": "secret",
				"fullname": "victor onah",
				"email":    "onahvictor@gmail.com"},
			buildStubs: func(repo *mockdb.MockUserRepository) {
				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(nil, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, r.Code)
			},
		},{
			name: "Wrong Email",
			body: map[string]string{
				"username": "victor",
				"password": "secret",
				"fullname": "victor onah",
				"email":    "onahvictorgmail.com"},
			buildStubs: func(repo *mockdb.MockUserRepository) {
				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, r.Code)
			},
		},
	}

	for _, value := range testCases {
		t.Run(value.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mockdb.NewMockUserRepository(ctrl)

			userSvc := service.NewUserService(repo)
			httpHandler := httptransport.NewUserHandler(userSvc)
			router := httptransport.NewRouter(nil, nil, httpHandler)

			data, err := json.Marshal(value.body)
			require.NoError(t, err)

			reader := bytes.NewReader([]byte(data))
			req, err := http.NewRequest(http.MethodPost, "/user", reader)
			require.NoError(t, err)

			value.buildStubs(repo)
			rec := httptest.NewRecorder()
			router.Mux.ServeHTTP(rec, req)
			value.checkResponse(t, rec)
		})
	}
}
