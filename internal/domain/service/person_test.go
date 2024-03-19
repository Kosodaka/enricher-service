package service

import (
	"context"
	"github.com/Kosodaka/enricher-service/internal/domain/dto"
	domainErr "github.com/Kosodaka/enricher-service/internal/domain/errors"
	"github.com/Kosodaka/enricher-service/internal/domain/model"
	"github.com/Kosodaka/enricher-service/internal/domain/ports/enricher"
	"github.com/Kosodaka/enricher-service/pkg/logger"
	mock_enricher "github.com/Kosodaka/enricher-service/pkg/mocks/api/enricher"
	mock_repository "github.com/Kosodaka/enricher-service/pkg/mocks/api/repository"
	mock_service "github.com/Kosodaka/enricher-service/pkg/mocks/api/service"
	"github.com/Kosodaka/enricher-service/pkg/validator"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
	"time"
)

type dependencies struct {
	repository *mock_repository.MockPersonRepository
	enricher   *mock_enricher.MockEnricher
	service    *mock_service.MockPersonService
}

func TestService_AddPerson(t *testing.T) {
	cases := []struct {
		name        string
		deadline    time.Duration
		input       *dto.AddPersonDTO
		enrichData  *enricher.EnrichData
		preparation func(d *dependencies, data *dto.AddPersonDTO, enrichData *enricher.EnrichData, ctx context.Context, err error)
		output      int
		err         error
	}{
		{
			name:     "valid input data",
			deadline: time.Second * 10,
			input: &dto.AddPersonDTO{
				Name:    "Oleg",
				Surname: "Dementiev",
			},
			enrichData: &enricher.EnrichData{
				Age:         60,
				Gender:      "male",
				Nationality: "RU",
			},
			preparation: func(d *dependencies, data *dto.AddPersonDTO, enrichData *enricher.EnrichData, ctx context.Context, err error) {
				person := &model.Person{
					Name:        data.Name,
					Surname:     data.Surname,
					Patronymic:  data.Patronymic,
					Age:         enrichData.Age,
					Gender:      enrichData.Gender,
					Nationality: enrichData.Nationality,
				}
				d.repository.EXPECT().AddPerson(ctx, person).Return(int(1), nil)
				d.enricher.EXPECT().Enrich(ctx, data.Name).Return(enrichData, nil)
			},
			output: 1,
			err:    nil,
		}, {
			name:     "invalid input data",
			deadline: time.Second * 10,
			input: &dto.AddPersonDTO{
				Name:    "danila",
				Surname: "Ivashenko",
			},
			enrichData:  nil,
			preparation: nil,
			output:      0,
			err:         validation.Errors{"name": domainErr.InvalidData},
		},
	}
	ctrl := gomock.NewController(t)
	valid := validator.NewValidator()
	dependencies := &dependencies{
		repository: mock_repository.NewMockPersonRepository(ctrl),
		enricher:   mock_enricher.NewMockEnricher(ctrl),
	}
	logger := logger.SetupLogger("test")
	svc := NewService()
	svc.Init(SetLogger(logger), SetValidator(valid))
	for _, testCases := range cases {
		t.Run(testCases.name, func(t *testing.T) {
			ctx, cansel := context.WithTimeout(context.Background(), testCases.deadline)
			defer cansel()
			if testCases.preparation != nil {
				testCases.preparation(dependencies, testCases.input, testCases.enrichData, ctx, testCases.err)
			}
			svc.opts.Repository = dependencies.repository
			svc.opts.Enricher = dependencies.enricher

			result, err := svc.AddPerson(ctx, testCases.input)
			if result != testCases.output {
				t.Errorf("got %d, want %d", result, testCases.output)
			}
			if !reflect.DeepEqual(err, testCases.err) {
				t.Errorf("got %v, want %v", err, testCases.err)
			}
		})
	}
}
func TestService_GetPerson(t *testing.T) {
	cases := []struct {
		name        string
		deadline    time.Duration
		input       int
		preparation func(d *dependencies, id int, ctx context.Context, err error)
		output      *model.Person
		err         error
	}{
		{
			name:     "valid input data",
			deadline: time.Second * 10,
			input:    1,
			preparation: func(d *dependencies, id int, ctx context.Context, err error) {
				person := &model.Person{
					Id:          1,
					Name:        "Nikolay",
					Surname:     "Chernyaev",
					Patronymic:  "",
					Age:         21,
					Gender:      "male",
					Nationality: "RU",
				}
				d.repository.EXPECT().GetPerson(ctx, id).Return(person, nil)
			},
			output: &model.Person{
				Id:          1,
				Name:        "Nikolay",
				Surname:     "Chernyaev",
				Patronymic:  "",
				Age:         21,
				Gender:      "male",
				Nationality: "RU",
			},
			err: nil,
		},
		{
			name:     "invalid id",
			deadline: time.Second * 10,
			input:    0,
			output:   nil,
			err:      domainErr.InvalidId,
		},
	}

	ctrl := gomock.NewController(t)
	valid := validator.NewValidator()
	dependencies := &dependencies{
		repository: mock_repository.NewMockPersonRepository(ctrl),
		enricher:   mock_enricher.NewMockEnricher(ctrl),
	}
	logger := logger.SetupLogger("test")
	svc := NewService()
	svc.Init(SetLogger(logger), SetValidator(valid))

	for _, testCases := range cases {
		t.Run(testCases.name, func(t *testing.T) {
			ctx, cansel := context.WithTimeout(context.Background(), testCases.deadline)
			defer cansel()
			if testCases.preparation != nil {
				testCases.preparation(dependencies, testCases.input, ctx, testCases.err)
			}
			svc.opts.Repository = dependencies.repository
			svc.opts.Enricher = dependencies.enricher
			result, err := svc.GetPerson(ctx, testCases.input)

			if !reflect.DeepEqual(result, testCases.output) {
				t.Errorf("got %v, want %v", result, testCases.output)
			}
			if !reflect.DeepEqual(err, testCases.err) {
				t.Errorf("got %v, want %v", err, testCases.err)
			}
		})

	}
}
