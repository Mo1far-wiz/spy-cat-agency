package main

import (
	"log"
	"net/http"
	"spy-cat-agency/internal/api"
	"spy-cat-agency/internal/db"
	"spy-cat-agency/internal/store"
	"time"

	"github.com/gin-gonic/gin"
)

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

func main() {
	cfg := Config{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
		DB: DBConfig{
			Addr:         "postgres://admin:adminpassword@localhost/agency?sslmode=disable",
			MaxOpenConns: 30,
			MaxIdleConns: 30,
			MaxIdleTime:  "15m",
		},
	}

	db, err := db.New(
		cfg.DB.Addr,
		cfg.DB.MaxOpenConns,
		cfg.DB.MaxIdleConns,
		cfg.DB.MaxIdleTime,
	)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	store := store.NewStorage(db)

	router := gin.Default()
	api.Mount(router)

	app := &Application{
		Config: cfg,
		Store:  store,
		Router: router,
	}
	app.run()
}

func (app *Application) run() {
	server := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      app.Router,
		ReadTimeout:  app.Config.ReadTimeout,
		WriteTimeout: app.Config.WriteTimeout,
		IdleTimeout:  app.Config.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
