package main

import (
	"github.com/Kosodaka/enricher-service/internal/adapters/app/app"
	"github.com/Kosodaka/enricher-service/internal/adapters/app/router"
	"github.com/Kosodaka/enricher-service/internal/adapters/enricher"
	"github.com/Kosodaka/enricher-service/internal/adapters/repository"
	"github.com/Kosodaka/enricher-service/internal/adapters/repository/postgres"
	"github.com/Kosodaka/enricher-service/internal/domain/service"
	"github.com/Kosodaka/enricher-service/pkg/config"
	"github.com/Kosodaka/enricher-service/pkg/logger"
	"github.com/Kosodaka/enricher-service/pkg/validator"
	"log/slog"
)

func main() {
	if err := config.LoadEnv(".env"); err != nil {
		panic(err)
	}

	cfg := config.LoadConfig()

	logger := logger.SetupLogger(cfg.GetEnv())

	logger.Info("start", slog.String("env", cfg.Env))
	enricher := enricher.NewEnricher(cfg)
	psql := postgres.NewPsql(cfg.PostgresDSN)
	db, err := psql.GetDb()
	if err != nil {
		panic(err)
	}

	personRepository := repository.NewPersonPostgres(db)
	personService := service.NewService(personRepository, enricher, logger, validator.NewValidator())
	personRouter := app.NewPersonRouter(personService)
	app := router.NewRouter(cfg, personRouter)
	if err := app.Run(); err != nil {
		panic(err)
	}

}
