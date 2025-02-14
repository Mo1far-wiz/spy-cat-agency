package main

import (
	"log"
	"spy-cat-agency/internal/api"
	"spy-cat-agency/internal/application"
	"spy-cat-agency/internal/db"
	"spy-cat-agency/internal/env"
	"spy-cat-agency/internal/store"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := application.Config{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
		DB: application.DBConfig{
			Addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/agency?sslmode=disable"),
			MaxOpenConns: env.GetInt("MAX_OPEN_CONNS", 30),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
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

	application.App = application.Application{
		Config: cfg,
		Store:  store,
		Router: router,
	}
	application.App.Run()
}
