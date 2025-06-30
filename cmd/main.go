package main

import (
	"log"

	"github.com/onahvictor/bank/config"
	"github.com/onahvictor/bank/db/client"
	"github.com/onahvictor/bank/db/repo"
	"github.com/onahvictor/bank/db/sqlc"
	"github.com/onahvictor/bank/service"
	httptransport "github.com/onahvictor/bank/transport/http"
)

func Run() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("failed to load config")
	}
	db, err := client.NewDBClient(config.DSN)
	if err != nil {
		log.Fatal("failed to connect to the datase", err)
	}

	store := sqlc.NewStore(db.Client)

	//repo setup
	accountRepo := repo.NewAccountRepo(store)
	transfRepo := repo.NewTransferRepo(store)

	//services setup
	accountSvc := service.NewAccountService(accountRepo)
	transferSvc := service.NewTransferService(transfRepo, accountRepo)
	//handlers
	accountHand := httptransport.NewAccountHandler(accountSvc)
	transfHand := httptransport.NewTranserHandler(transferSvc)

	//router & routes setup
	router := httptransport.NewRouter(accountHand,transfHand)

	if err := router.Serve(config.PORT); err != nil {
		log.Fatal(err)
	}
}

func main() {
	Run()
}
