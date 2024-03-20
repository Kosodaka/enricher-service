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
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const billingAddr = "http://localhost:8080"

type TestSuite struct {
	suite.Suite
	psqlContainer *PostgreSQLContainer
	server        *httptest.Server
}

func (s *TestSuite) SetupSuite() {
	// create db container
	ctx, ctxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer ctxCancel()

	psqlContainer, err := NewPostgreSQLContainer(ctx)
	s.Require().NoError(err)

	s.psqlContainer = psqlContainer

	//err = migrate.Migrate(psqlContainer.GetDSN(), migrate.Migrations)
	s.Require().NoError(err)

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

func (s *TestSuite) TestAddPerson(t *testing.T) {
	id := 1
	Test(t,
		Post("localhost:8080/persons"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().JSON(&dto.AddPersonDTO{
			Name:       "Aleks",
			Surname:    "Ivanov",
			Patronymic: "",
		}),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().JSON().JQ(".json.id").Equal(id),
	)
}
func TestSuite_Run(t *testing.T) {
	suite.Run(t, new(TestSuite))

}

func (s *TestSuite) TearDownSuite() {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()

	s.Require().NoError(s.psqlContainer.Terminate(ctx))
	s.server.Close()

}

func (s *TestSuite) TestGetPerson(t *testing.T) {
	Test(t,
		Description("Get to localhost:8080"),
		Get("http://localhost:8080/person/:1"),
		Expect().Status().Equal(http.StatusOK),
	)

}
