package service

import (
	"context"
	"github.com/Kosodaka/enricher-service/internal/domain/model"
	"github.com/Kosodaka/enricher-service/internal/domain/ports/enricher"
	"github.com/Kosodaka/enricher-service/internal/domain/ports/repository"
	"log/slog"
)

type Validator interface {
	ValidateId(id int) error
	ValidateDataToAdd(data *PersonFullName) error
	ValidateDataToGet(data *model.Person) error
	ValidateDataToUpdate(data *model.Person) error
}

type PersonFullName struct {
	Id         int    `json:"id,string"`
	Name       string `json:"name" db:"name"`
	Surname    string `json:"surname" db:"surname"`
	Patronymic string `json:"patronymic" db:"patronymic"`
}
type Service struct {
	Repository repository.PersonRepository
	Enricher   enricher.Enricher
	Logger     *slog.Logger
	Validator  Validator
}

func NewService(r repository.PersonRepository, e enricher.Enricher, l *slog.Logger, v Validator) *Service {
	return &Service{
		Repository: r,
		Enricher:   e,
		Logger:     l,
		Validator:  v,
	}
}

func (s Service) AddPerson(ctx context.Context, data *PersonFullName) (int, error) {
	op := "service.AddPerson"
	logger := s.Logger.With("operation", op)
	if err := s.Validator.ValidateDataToAdd(data); err != nil {
		return 0, err
	}
	enrichData, err := s.Enricher.Enrich(ctx, data.Name)
	if err != nil {
		return 0, err
	}

	personModel := &model.Person{
		Name:        data.Name,
		Surname:     data.Surname,
		Patronymic:  data.Patronymic,
		Age:         enrichData.Age,
		Gender:      enrichData.Gender,
		Nationality: enrichData.Nationality,
	}

	id, err := s.Repository.AddPerson(ctx, personModel)
	if err != nil {
		logger.Debug("failed to add person", slog.Any("error", err))
		return 0, err
	}

	logger.Debug("person was successfully added", slog.Any("id", id))
	return id, nil
}

func (s Service) GetPerson(ctx context.Context, id int) (*model.Person, error) {
	op := "service.GetPerson"
	logger := s.Logger.With("operation", op)
	if err := s.Validator.ValidateId(id); err != nil {
		return nil, err
	}
	person, err := s.Repository.GetPerson(ctx, id)
	if err != nil {

	}
	logger.Debug("person was successfully got", slog.Int("id", id))
	return person, nil
}

func (s Service) GetPersons(ctx context.Context, data *model.Person) ([]model.Person, error) {
	op := "service.GetPersons"
	logger := s.Logger.With("operation", op)
	if err := s.Validator.ValidateDataToGet(data); err != nil {
		return nil, err
	}
	person, err := s.Repository.GetPersons(ctx, data)
	if err != nil {
		logger.Debug("failed to get persons")
		return nil, err
	}
	logger.Debug("persons was successfully got")
	return person, nil
}

func (s Service) DeletePerson(ctx context.Context, id int) error {
	op := "service.DeletePerson"
	logger := s.Logger.With("operation", op)
	if err := s.Validator.ValidateId(id); err != nil {
		return err
	}
	err := s.Repository.DeletePerson(ctx, id)
	if err != nil {
		logger.Debug("fail to delete person", slog.Int("id", id), slog.String("error", err.Error()))
	} else {
		logger.Debug("person was successfully deleted", slog.Int("id", id))
	}
	return err
}

func (s Service) UpdatePerson(ctx context.Context, data *model.Person) error {
	op := "service.DeletePerson"
	logger := s.Logger.With("operation", op)
	if err := s.Validator.ValidateDataToUpdate(data); err != nil {
		return err
	}
	err := s.Repository.UpdatePerson(ctx, data)
	if err != nil {
		logger.Debug("fail to update person", slog.Any("error", err.Error()))
	} else {
		logger.Debug("person was successfully updated")
	}
	return err
}
