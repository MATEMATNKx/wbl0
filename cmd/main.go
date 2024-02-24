package main

import (
	"context"
	"errors"
	"l0/internal/config"
	"l0/internal/service/cache"
	"l0/internal/service/order"
	"l0/internal/service/stan"
	"l0/internal/service/storage"
	"l0/internal/transport"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.New()
	orderSvc := order.New(
		storage.New(cfg.DB),
		cache.New(),
	)

	stan := stan.New(cfg.Stan, orderSvc)

	app := transport.New(cfg, orderSvc, stan)
	go func() {
		if err := app.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error on start web server: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Stop(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
