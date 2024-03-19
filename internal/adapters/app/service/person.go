package service

import (
	"context"
	"github.com/Kosodaka/enricher-service/internal/domain/dto"
	"github.com/Kosodaka/enricher-service/internal/domain/model"
)

type PersonService interface {
	AddPerson(ctx context.Context, data *dto.AddPersonDTO) (int, error)
	GetPerson(ctx context.Context, id int) (*model.Person, error)
	GetPersons(ctx context.Context, data *model.Person) ([]model.Person, error)
	UpdatePerson(ctx context.Context, data *model.Person) error
	DeletePerson(ctx context.Context, id int) error
}
