package httptransport

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Router struct {
	Mux *gin.Engine
}

func NewRouter(accountHand *AccountHandler, transferHand *TransferHandler) *Router {
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	accountHand.MapAccountRoutes(router)
	transferHand.MapAccountRoutes(router)

	routerSetup := &Router{
		Mux: router,
	}
	return routerSetup
}

func (r *Router) Serve(port string) error {
	server := &http.Server{
		Addr:         port,
		Handler:      r.Mux,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 15,
	}

	shutDown := make(chan error)
	go func() {
		err := server.ListenAndServe()
		shutDown <- err
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	fmt.Printf("%s, signal recieved shutting down server\n", (<-quit).String())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	err := server.Shutdown(ctx) //is htis a blocking till error returns
	if err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	err = <-shutDown
	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server closed unexpectedly: %w", err)
	}

	fmt.Println("Server shutdown complete")
	return nil
}
