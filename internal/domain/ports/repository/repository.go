package repository

import (
	"context"
	"github.com/Kosodaka/enricher-service/internal/domain/model"
)

type PersonRepository interface {
	AddPerson(context.Context, *model.Person) (int, error)
	GetPerson(context.Context, int) (*model.Person, error)
	GetPersons(context.Context, *model.Person) ([]model.Person, error)
	UpdatePerson(context.Context, *model.Person) error
	DeletePerson(context.Context, int) error
}
