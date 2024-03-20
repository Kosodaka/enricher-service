//go:build integration
// +build integration

package integration

import (
	"context"
	. "github.com/Eun/go-hit"
	"github.com/Kosodaka/enricher-service/internal/adapters/app/app"
	"github.com/Kosodaka/enricher-service/internal/adapters/app/router"
	"github.com/Kosodaka/enricher-service/internal/adapters/enricher"
	"github.com/Kosodaka/enricher-service/internal/adapters/repository"
	"github.com/Kosodaka/enricher-service/internal/adapters/repository/postgres"
	"github.com/Kosodaka/enricher-service/internal/domain/dto"
	"github.com/Kosodaka/enricher-service/internal/domain/service"
	"github.com/Kosodaka/enricher-service/pkg/config"
	"github.com/Kosodaka/enricher-service/pkg/logger"
	"github.com/Kosodaka/enricher-service/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

type TestSuite struct {
	suite.Suite
	psqlContainer *PostgreSQLContainer
	server        *gin.Engine
}

func (s *TestSuite) SetupSuite() {
	// create db container
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	psqlContainer, err := NewPostgreSQLContainer(ctx)
	s.Require().NoError(err)

	s.psqlContainer = psqlContainer

	//main
	if err := config.LoadEnv(".env"); err != nil {
		panic(err)
	}
	cfg := config.LoadConfig()
	logger := logger.SetupLogger(cfg.GetEnv())
	valid := validator.NewValidator()
	psql := postgres.NewPsql(psqlContainer.GetDSN())
	db, err := psql.GetDb()
	if err != nil {
		panic(err)
	}
	enricher := enricher.NewEnricher(cfg)
	personRepository := repository.NewPersonPostgres(db)
	personService := service.NewService()
	personService.Init(service.SetRepository(personRepository), service.SetEnricher(enricher), service.SetLogger(logger), service.SetValidator(valid))
	personRouter := app.NewPersonRouter(personService)
	app := router.NewRouter(cfg, personRouter)
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func TestSuite_Run(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
func (s *TestSuite) TestAddPerson(t *testing.T) {
	id := 1
	Test(t,
		Post("localhost:8080/persons"),
		Send().Headers("Content-Type").Add("application/json"),
		Debug(),
		Send().Body().JSON(&dto.AddPersonDTO{
			Name:       "Aleks",
			Surname:    "Ivanov",
			Patronymic: "",
		}),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().JSON().JQ(".json.id").Equal(id),
	)
}

func (s *TestSuite) TestGetPerson(t *testing.T) {
	Test(t,
		Description("Get to localhost:8080"),
		Get("http://localhost:8080/person/:1"),
		Expect().Status().Equal(http.StatusOK),
	)

}
func (s *TestSuite) TearDownSuite() {
	s.server = gin.Default()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: s.server,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
