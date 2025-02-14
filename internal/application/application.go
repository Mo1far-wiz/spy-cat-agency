package application

import (
	"net/http"
	"spy-cat-agency/internal/store"
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

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
