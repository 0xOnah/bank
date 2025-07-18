package main

import (
	"log"

	"github.com/onahvictor/bank/internal/config"
	"github.com/onahvictor/bank/internal/db/client"
	"github.com/onahvictor/bank/internal/db/repo"
	"github.com/onahvictor/bank/internal/db/sqlc"
	"github.com/onahvictor/bank/internal/sdk/auth"
	"github.com/onahvictor/bank/internal/service"
	httptransport "github.com/onahvictor/bank/internal/transport/http"
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

	auth, err := auth.NewJWTMaker(config.TOKEN_SYMMETRIC_KEY)
	if err != nil {
		log.Fatal("auth token creation failed", err)
	}
	//repo setup
	accountRepo := repo.NewAccountRepo(store)
	transfRepo := repo.NewTransferRepo(store)
	userRepo := repo.NewUserRepo(store)
	sessionRepo := repo.NewSessionRepo(store)

	//services setup
	accountSvc := service.NewAccountService(accountRepo)
	transferSvc := service.NewTransferService(transfRepo, accountRepo)
	usrSvc := service.NewUserService(userRepo, auth, config, sessionRepo)
	//handlers

	accountHand := httptransport.NewAccountHandler(accountSvc, auth)
	transfHand := httptransport.NewTranserHandler(transferSvc, auth)
	userHand := httptransport.NewUserHandler(usrSvc, auth)

	//router & routes setup
	router := httptransport.NewRouter(accountHand, transfHand, userHand)

	if err := router.Serve(config.PORT); err != nil {
		log.Fatal(err)
	}
}

func main() {
	Run()
}
