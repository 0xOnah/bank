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

	mockdb "github.com/onahvictor/bank/internal/db/mock"
	"github.com/onahvictor/bank/internal/db/repo"
	"github.com/onahvictor/bank/internal/entity"
	"github.com/onahvictor/bank/internal/service"
	httptransport "github.com/onahvictor/bank/internal/transport/http"
	"github.com/onahvictor/bank/internal/util"
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
	require.NoError(t, err)
	require.Equal(t, account.ID, got.ID)
	require.Equal(t, account.Balance, got.Balance)
	require.Equal(t, account.Currency, got.Currency)

}

func TestGetAccountByID(t *testing.T) {
	expected := randomAccount()

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
			accountSvc := service.NewAccountService(accountRepo)
			httpHandler := httptransport.NewAccountHandler(accountSvc)
			router := httptransport.NewRouter(httpHandler, nil,nil)

			value.buildStubs(accountRepo)

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts/%d", value.accountID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			router.Mux.ServeHTTP(recorder, req)

			value.checkResponse(t, recorder)
		})
	}

}
