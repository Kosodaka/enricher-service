package service

import (
	"context"
	"github.com/Kosodaka/enricher-service/internal/domain/dto"
	"github.com/Kosodaka/enricher-service/internal/domain/model"
	"github.com/Kosodaka/enricher-service/internal/domain/ports/enricher"
	"github.com/Kosodaka/enricher-service/internal/domain/ports/repository"
	"log/slog"
)

type Validator interface {
	ValidateId(id int) error
	ValidateDataToAdd(data *dto.AddPersonDTO) error
	ValidateDataToGet(data *model.Person) error
	ValidateDataToUpdate(data *model.Person) error
}

type PersonFullName struct {
	Id         int    `json:"id,string"`
	Name       string `json:"name" db:"name"`
	Surname    string `json:"surname" db:"surname"`
	Patronymic string `json:"patronymic" db:"patronymic"`
}

type Service interface {
	Repository() repository.PersonRepository
	Enricher() enricher.Enricher
	Logger() slog.Logger
	Validator() Validator
	Init(...Option)

	AddPerson(ctx context.Context, data *dto.AddPersonDTO) (int, error)
	GetPerson(ctx context.Context, id int) (*model.Person, error)
	GetPersons(ctx context.Context, data *model.Person) ([]model.Person, error)
	UpdatePerson(ctx context.Context, data *model.Person) error
	DeletePerson(ctx context.Context, id int) error
}
type Options struct {
	Repository repository.PersonRepository
	Enricher   enricher.Enricher
	Logger     *slog.Logger
	Validator  Validator
}

type Option func(*Options) error

func NewOptions(opts ...Option) Options {
	options := Options{
		Repository: repository.PersonRepository(nil),
		Enricher:   enricher.Enricher(nil),
		Logger:     &slog.Logger{},
		Validator:  Validator(nil),
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

func NewService(opts ...Option) *service {
	return &service{
		opts: NewOptions(opts...),
	}
}

type service struct {
	opts Options
}

func (s *service) Init(opts ...Option) error {
	var err error
	// process options
	for _, o := range opts {
		if err = o(&s.opts); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) Repository() repository.PersonRepository {
	return s.opts.Repository
}

func (s *service) Enricher() enricher.Enricher {
	return s.opts.Enricher
}

func (s *service) Logger() *slog.Logger {
	return s.opts.Logger
}

func (s *service) Validator() Validator {
	return s.opts.Validator
}

func SetRepository(r repository.PersonRepository) Option {
	return func(o *Options) error {
		o.Repository = r
		return nil
	}
}

func SetEnricher(e enricher.Enricher) Option {
	return func(o *Options) error {
		o.Enricher = e
		return nil
	}
}

func SetLogger(l *slog.Logger) Option {
	return func(o *Options) error {
		o.Logger = l
		return nil
	}
}

func SetValidator(v Validator) Option {
	return func(o *Options) error {
		o.Validator = v
		return nil
	}
}

func (s service) AddPerson(ctx context.Context, data *dto.AddPersonDTO) (int, error) {
	op := "service.AddPerson"
	logger := s.opts.Logger.With("operation", op)
	if err := s.opts.Validator.ValidateDataToAdd(data); err != nil {
		return 0, err
	}
	enrichData, err := s.opts.Enricher.Enrich(ctx, data.Name)
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

	id, err := s.opts.Repository.AddPerson(ctx, personModel)
	if err != nil {
		logger.Debug("failed to add person", slog.Any("error", err))
		return 0, err
	}

	logger.Debug("person was successfully added", slog.Any("id", id))
	return id, nil
}

func (s service) GetPerson(ctx context.Context, id int) (*model.Person, error) {
	op := "service.GetPerson"
	logger := s.opts.Logger.With("operation", op)
	if err := s.opts.Validator.ValidateId(id); err != nil {
		return nil, err
	}
	person, err := s.opts.Repository.GetPerson(ctx, id)
	if err != nil {
		return nil, err
	}
	logger.Debug("person was successfully got", slog.Int("id", id))
	return person, nil
}

func (s service) GetPersons(ctx context.Context, data *model.Person) ([]model.Person, error) {
	op := "service.GetPersons"
	logger := s.opts.Logger.With("operation", op)
	if err := s.opts.Validator.ValidateDataToGet(data); err != nil {
		return nil, err
	}
	person, err := s.opts.Repository.GetPersons(ctx, data)
	if err != nil {
		logger.Debug("failed to get persons")
		return nil, err
	}
	logger.Debug("persons was successfully got")
	return person, nil
}

func (s service) DeletePerson(ctx context.Context, id int) error {
	op := "service.DeletePerson"
	logger := s.opts.Logger.With("operation", op)
	if err := s.opts.Validator.ValidateId(id); err != nil {
		return err
	}
	err := s.opts.Repository.DeletePerson(ctx, id)
	if err != nil {
		logger.Debug("fail to delete person", slog.Int("id", id), slog.String("error", err.Error()))
	} else {
		logger.Debug("person was successfully deleted", slog.Int("id", id))
	}
	return err
}

func (s service) UpdatePerson(ctx context.Context, data *model.Person) error {
	op := "service.DeletePerson"
	logger := s.opts.Logger.With("operation", op)
	if err := s.opts.Validator.ValidateDataToUpdate(data); err != nil {
		return err
	}
	err := s.opts.Repository.UpdatePerson(ctx, data)
	if err != nil {
		logger.Debug("fail to update person", slog.Any("error", err.Error()))
	} else {
		logger.Debug("person was successfully updated")
	}
	return err
}
