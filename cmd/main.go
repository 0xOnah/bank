package main

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/0xOnah/bank/doc"
	"github.com/0xOnah/bank/internal/config"
	"github.com/0xOnah/bank/internal/db"
	"github.com/0xOnah/bank/internal/db/repo"
	"github.com/0xOnah/bank/internal/db/sqlc"
	"github.com/0xOnah/bank/internal/sdk/auth"
	"github.com/0xOnah/bank/internal/sdk/jobs"
	"github.com/0xOnah/bank/internal/sdk/logger"
	"github.com/0xOnah/bank/internal/service"
	grpctransport "github.com/0xOnah/bank/internal/transport/grpc"
	httptransport "github.com/0xOnah/bank/internal/transport/http"
	"github.com/0xOnah/bank/internal/transport/sdk/middleware"
	"github.com/0xOnah/bank/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	//config
	config, err := config.LoadConfig(".")
	if err != nil {
		return
	}

	logger, err := logger.InitLogger(&config)
	if err != nil {
		// using global log to kill program if our custom logger is not initizalied
		log.Fatal().Err(err).Send()
	}

	//db setup
	logger.Info().Msg(config.DSN)
	database, err := db.NewDBClient(config.DSN)
	if err != nil {
		logger.Fatal().Err(err).Msg("invalid database.DSN")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = database.Ping(ctx) //ensure connection
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to the database")
	}
	logger.Info().Msg("database connection esterblished")

	err = database.MigrateUP() //migration
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to run up-migration")
	}
	logger.Info().Msg("database migration successful")

	store := sqlc.NewStore(database.Client)

	//redis
	redisOpts := asynq.RedisClientOpt{
		Addr: config.REDIS_ADDRESS,
	}
	taskQueue := jobs.NewTaskQueue(redisOpts, logger)

	//authenticator
	auth, err := auth.NewJWTMaker(config.TOKEN_SYMMETRIC_KEY)
	if err != nil {
		logger.Fatal().Err(err).Msg("jwt-maker not initialized")
	}
	// RunHttpServer(store, config, auth)
	go runJobService(redisOpts, store, logger)
	RunGrpcServer(config, store, auth, logger, taskQueue)
	// RunGatewayServer(config, store, auth, logger)
}

func runJobService(redisOpts asynq.RedisClientOpt, store *sqlc.SQLStore, logger *zerolog.Logger) {
	UserRepo := repo.NewUserRepo(store)
	taskProcessor := jobs.NewWorkerService(redisOpts, UserRepo, logger)
	log.Info().Msg("starting task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create redis server")
	}
}
func RunHttpServer(
	store *sqlc.SQLStore,
	config config.Config,
	auth auth.Authenticator,
) {
	accountRepo := repo.NewAccountRepo(store)
	transfRepo := repo.NewTransferRepo(store)
	UserRepo := repo.NewUserRepo(store)
	sessionRepo := repo.NewSessionRepo(store)

	//services setup
	accountSvc := service.NewAccountService(accountRepo)
	transferSvc := service.NewTransferService(transfRepo, accountRepo)
	usrSvc := service.NewUserService(UserRepo, auth, config, sessionRepo)
	//handlers
	accountHand := httptransport.NewAccountHandler(accountSvc, auth)
	transfHand := httptransport.NewTranserHandler(transferSvc, auth)
	userHand := httptransport.NewUserHandler(usrSvc, auth)

	//router & routes setup
	router := httptransport.NewRouter(accountHand, transfHand, userHand)

	if err := router.Serve(config.HTTP_SERVER_ADDRESS); err != nil {
		return
	}
}

func RunGatewayServer(
	config config.Config,
	store *sqlc.SQLStore,
	tokenMaker auth.Authenticator,
	log *zerolog.Logger,
	taskqueue jobs.TaskDistributor,

) {
	ur := repo.NewUserRepo(store)
	sr := repo.NewSessionRepo(store)
	UserRepo := repo.NewUserRepo(store)

	usrSvc := service.NewUserService(ur, tokenMaker, config, sr)
	svcLogger := logger.ServiceLogger(log, "auth_Service")
	UserHandler := grpctransport.NewUserHandler(usrSvc, UserRepo, tokenMaker, svcLogger, taskqueue)

	httpGateWayMux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := pb.RegisterUserServiceHandlerServer(ctx, httpGateWayMux, UserHandler)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register userHandler with the server")
	}

	httpmux := http.NewServeMux()
	httpmux.Handle("/", httpGateWayMux)

	fs := http.FileServer(http.FS(doc.SwaggerFs))
	httpmux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	httpmux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
			w.WriteHeader(500)
		}
	})

	listener, err := net.Listen("tcp", config.HTTP_SERVER_ADDRESS)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create grpc-gateway listener")
	}

	log.Info().Str("port", config.GRPC_SERVER_ADDRESS).Msg("starting grpc-gateway server")
	reqlog := middleware.LogRequest(log)
	err = http.Serve(listener, reqlog(httpmux))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start-up grpc-gateway server")
	}
}

// starting a grpc server
func RunGrpcServer(
	config config.Config,
	store *sqlc.SQLStore,
	tokenMaker auth.Authenticator,
	log *zerolog.Logger,
	taskqueue jobs.TaskDistributor,
) {
	ur := repo.NewUserRepo(store)
	sr := repo.NewSessionRepo(store)
	UserRepo := repo.NewUserRepo(store)
	usrSvc := service.NewUserService(ur, tokenMaker, config, sr)
	UserHandler := grpctransport.NewUserHandler(usrSvc, UserRepo, tokenMaker, log, taskqueue)

	logger := grpctransport.LoggingInterceptor(log)
	recoverPanic := grpctransport.UnaryRecoverPanicInterceptor(log)

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(recoverPanic, logger))

	reflection.Register(grpcServer)

	pb.RegisterUserServiceServer(grpcServer, UserHandler)

	listener, err := net.Listen("tcp", config.GRPC_SERVER_ADDRESS)
	if err != nil {
		log.Error().Err(err).Msg("failed to create grpc listener")
	}

	log.Info().Str("port", config.GRPC_SERVER_ADDRESS).Msg("starting grpc server")
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Error().Err(err).Msg("failed to statup grpc server")
	}
}
