package application

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"spy-cat-agency/internal/store"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

var App Application

type Application struct {
	Config Config
	Store  store.Storage
	Router *gin.Engine
}

type Config struct {
	Addr         string
	DB           DBConfig
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	IdleTimeout  time.Duration
}

type DBConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

func (app *Application) Run() {
	server := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      app.Router,
		ReadTimeout:  app.Config.ReadTimeout,
		WriteTimeout: app.Config.WriteTimeout,
		IdleTimeout:  app.Config.IdleTimeout,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		log.Println("Server is starting on", app.Config.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server gracefully stopped")

	wg.Wait()
}
