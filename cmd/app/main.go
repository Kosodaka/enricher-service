package main

import (
	"fmt"
	"github.com/Kosodaka/enricher-service/internal/adapters/app/app"
	"github.com/Kosodaka/enricher-service/internal/adapters/app/router"
	"github.com/Kosodaka/enricher-service/internal/adapters/enricher"
	"github.com/Kosodaka/enricher-service/internal/adapters/repository"
	"github.com/Kosodaka/enricher-service/internal/adapters/repository/postgres"
	"github.com/Kosodaka/enricher-service/internal/domain/service"
	"github.com/Kosodaka/enricher-service/pkg/config"
	"github.com/Kosodaka/enricher-service/pkg/logger"
	"log"
	"log/slog"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("error with loading config", err)
	}
	fmt.Println(cfg)
	logger := logger.SetupLogger(cfg.Env)

	logger.Info("start", slog.String("env", cfg.Env))
	enricher := enricher.NewEnricher(cfg)
	psql := postgres.NewPsql(cfg.PostgresDSN)
	db, err := psql.GetDb()
	if err != nil {
		logger.Error("fail to init db", logger.LogAttrs)
	}

	personRepository := repository.NewPersonPostgres(db)
	personService := service.NewService(personRepository, enricher, logger)
	personRouter := app.NewPersonRouter(personService)
	app := router.NewRouter(cfg, personRouter)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
