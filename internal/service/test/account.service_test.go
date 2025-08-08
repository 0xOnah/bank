package service_test

//this test uses the mock for testing
import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/0xOnah/bank/internal/config"
	mockdb "github.com/0xOnah/bank/internal/db/mock"
	"github.com/0xOnah/bank/internal/db/repo"
	"github.com/0xOnah/bank/internal/entity"
	"github.com/0xOnah/bank/internal/sdk/auth"
	"github.com/0xOnah/bank/internal/sdk/util"
	"github.com/0xOnah/bank/internal/service"
	httptransport "github.com/0xOnah/bank/internal/transport/http"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func randomAccount() *entity.Account {
	return &entity.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatch(t *testing.T, body *bytes.Buffer, account *entity.Account) {
	t.Helper()

	var got entity.Account
	err := json.Unmarshal(body.Bytes(), &got)
	t.Log(got)
	require.NoError(t, err)
	require.Equal(t, account.ID, got.ID)
	require.Equal(t, account.Owner, got.Owner)
	require.Equal(t, account.Balance, got.Balance)
	require.Equal(t, account.Currency, got.Currency)

}

func TestGetAccountByID(t *testing.T) {
	token, err := auth.NewJWTMaker("123456789123456789123456789123456789")
	require.NoError(t, err)

	expected := randomAccount()
	payload, _, err := token.GenerateToken(expected.Owner, time.Minute*15)
	require.NoError(t, err)

	type TestCase struct {
		name          string
		accountID     int64
		buildStubs    func(*mockdb.MockAccountRepository)
		checkResponse func(*testing.T, *httptest.ResponseRecorder)
	}

	testCases := []TestCase{
		{
			name:      "Ok: Account Found",
			accountID: expected.ID,
			buildStubs: func(mar *mockdb.MockAccountRepository) {
				mar.
					EXPECT().
					GetAccountByID(gomock.Any(), gomock.Eq(expected.ID)).
					Times(1).
					Return(expected, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				requireBodyMatch(t, recorder.Body, expected)
				require.Equal(t, http.StatusOK, recorder.Code)

			},
		}, {
			name:      "Error: Not Found",
			accountID: expected.ID,
			buildStubs: func(mar *mockdb.MockAccountRepository) {
				mar.
					EXPECT().
					GetAccountByID(gomock.Any(), gomock.Eq(expected.ID)).
					Times(1).
					Return(nil, repo.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		}, {
			name:      "Error: Internal Server Error",
			accountID: expected.ID,
			buildStubs: func(mar *mockdb.MockAccountRepository) {
				mar.
					EXPECT().
					GetAccountByID(gomock.Any(), gomock.Eq(expected.ID)).
					Times(1).
					Return(nil, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		}, {
			name:      "Error: Validation Failed",
			accountID: 0,
			buildStubs: func(mar *mockdb.MockAccountRepository) {
				mar.
					EXPECT().
					GetAccountByID(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, value := range testCases {
		t.Run(value.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			accountRepo := mockdb.NewMockAccountRepository(ctrl)
			transferRepo := mockdb.NewMockTransferRepository(ctrl)
			UserRepo := mockdb.NewMockUserRepository(ctrl)
			sessionRepo := mockdb.NewMockSessionRepository(ctrl)
			accountSvc := service.NewAccountService(accountRepo)
			accountHandler := httptransport.NewAccountHandler(accountSvc, token)

			transferSvc := service.NewTransferService(transferRepo, accountRepo)
			transfHand := httptransport.NewTranserHandler(transferSvc, token)

			usrSvc := service.NewUserService(UserRepo, token, config.Config{}, sessionRepo)
			userHand := httptransport.NewUserHandler(usrSvc, token)

			router := httptransport.NewRouter(accountHandler, transfHand, userHand)

			value.buildStubs(accountRepo)

			recorder := httptest.NewRecorder()
			t.Log(value.accountID)
			url := fmt.Sprintf("/accounts/%d", value.accountID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", payload))
			require.NoError(t, err)

			router.Mux.ServeHTTP(recorder, req)
			value.checkResponse(t, recorder)
		})
	}

}
