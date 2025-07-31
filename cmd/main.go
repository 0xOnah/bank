package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/0xOnah/bank/doc"
	"github.com/0xOnah/bank/internal/config"
	"github.com/0xOnah/bank/internal/db/client"
	"github.com/0xOnah/bank/internal/db/repo"
	"github.com/0xOnah/bank/internal/db/sqlc"
	"github.com/0xOnah/bank/internal/sdk/auth"
	"github.com/0xOnah/bank/internal/service"
	grpctransport "github.com/0xOnah/bank/internal/transport/grpc"
	httptransport "github.com/0xOnah/bank/internal/transport/http"
	"github.com/0xOnah/bank/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}).WithAttrs([]slog.Attr{})

	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	config, err := config.LoadConfig(".")
	if err != nil {
		slog.Error("msg", slog.Any("failed to load config", err))
		os.Exit(1)
	}
	db, err := client.NewDBClient(config.DSN)
	if err != nil {
		slog.Error("msg", slog.Any("failed to connect to the datase", err))
		os.Exit(1)
	}

	store := sqlc.NewStore(db.Client)

	auth, err := auth.NewJWTMaker(config.TOKEN_SYMMETRIC_KEY)
	if err != nil {
		slog.Error("msg", slog.Any("auth token creation failed", err))
		os.Exit(1)
	}
	//repo setup
	// RunHttpServer(store, config, auth)
	// go func() { RunGrpcServer(config, store, auth) }()
	RunGatewayServer(config, store, auth)
}

func RunHttpServer(store *sqlc.SQLStore, config config.Config, auth auth.Authenticator) {
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

	if err := router.Serve(config.HTTP_SERVER_ADDRESS); err != nil {
		slog.Error("msg", slog.Any("failed to run http server", err))
		os.Exit(1)
	}
}

func RunGatewayServer(config config.Config, store *sqlc.SQLStore, tokenMaker auth.Authenticator) {
	ur := repo.NewUserRepo(store)
	sr := repo.NewSessionRepo(store)
	usrSvc := service.NewUserService(ur, tokenMaker, config, sr)
	UserHandler := grpctransport.NewUserHandler(usrSvc)

	grpcmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := pb.RegisterUserServiceHandlerServer(ctx, grpcmux, UserHandler)
	if err != nil {
		log.Fatal("failed to register the userServiceHandlerServer ")
	}

	httpmux := http.NewServeMux()
	httpmux.Handle("/", grpcmux)

	fs := http.FileServer(http.FS(doc.SwaggerFs))
	httpmux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	httpmux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"server is live"}`)
	})

	listener, err := net.Listen("tcp", config.HTTP_SERVER_ADDRESS)
	if err != nil {
		log.Fatal("failed to create listner", err)
	}

	log.Println("server started Port:", config.HTTP_SERVER_ADDRESS)
	err = http.Serve(listener, httpmux)
	if err != nil {
		log.Fatal("failed to startup server")
	}
}

// starting a grpc server
func RunGrpcServer(config config.Config, store *sqlc.SQLStore, tokenMaker auth.Authenticator) {
	ur := repo.NewUserRepo(store)
	sr := repo.NewSessionRepo(store)
	usrSvc := service.NewUserService(ur, tokenMaker, config, sr)

	grpcServer := grpc.NewServer()
	UserHandler := grpctransport.NewUserHandler(usrSvc)
	reflection.Register(grpcServer) //imagince it as a self docnumentation for the server
	pb.RegisterUserServiceServer(grpcServer, UserHandler)

	listener, err := net.Listen("tcp", config.GRPC_SERVER_ADDRESS)
	if err != nil {
		slog.Error("msg", slog.Any("cannot cannot create listerner", err))
		os.Exit(1)
	}

	log.Printf("starting grpc server at port %s", config.GRPC_SERVER_ADDRESS)
	err = grpcServer.Serve(listener)
	if err != nil {
		slog.Error("msg", slog.Any("cannot start grpc client", err))
		os.Exit(1)
	}
}

//
