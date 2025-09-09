package main

import (
	"context"
	"fmt"
	"log"
	"ltp/internal/config"
	"ltp/internal/repository"
	"ltp/internal/service"
	"ltp/internal/transport"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	conf, err := config.New()

	if err != nil {
		logger.Fatal("Server Config Broken")
		os.Exit(1)
	}

	repo := repository.New(conf.KrakenAPIConfig, logger)
	svc := service.New(repo, conf.KrakenAPIConfig, logger)
	rtr := transport.NewHTTPRouter(svc, conf, logger)

	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%s", conf.ServerPort),
		// Setting Timeouts to avoid slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		// Router contain handlers
		Handler: rtr,
	}

	//Run server in go routine to prevent blocking

	go func() {
		logger.Printf("Server starting on: %s", conf.ServerPort)
		if err := srv.ListenAndServe(); err != nil {
			logger.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logger.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
