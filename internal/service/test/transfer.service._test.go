package service_test

// const currency = "USD"

// func TestCreateTransferTX(t *testing.T) {
// 	account1 := randomAccount()
// 	account1.Currency = currency
// 	account2 := randomAccount()
// 	account2.Currency = currency
// 	amount := 30

// 	transferRepo := repo.NewTransferRepo(testStore)
// 	accountRepo := repo.NewAccountRepo(testStore)

// 	transferService := service.NewTransferService(transferRepo, accountRepo)

// 	for i := 0; i <= 10; i++ {
// 		go func() {
// 			result, err := transferService.CreateTransferTX(context.Background(), entity.CreateTransferInput{
// 				FromAccountID: account1.ID,
// 				ToAccountID:   account2.ID,
// 				Amount:        int64(amount),
// 			}, currency)
// 			require.NoError(t, err)
// 		}()

// 	}
// }
