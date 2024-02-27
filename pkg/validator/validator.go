package validator

import (
	"errors"
	"github.com/Kosodaka/enricher-service/internal/domain/model"
	"github.com/Kosodaka/enricher-service/internal/domain/service"
	validation "github.com/go-ozzo/ozzo-validation"
	"regexp"
)

type Validator struct {
}

func NewValidator() *Validator {
	return &Validator{}
}
func (Validator) ValidateId(id int) error {
	err := errors.New("invalid id")
	if id <= 0 {
		return err
	}
	return nil
}

func (Validator) ValidateDataToAdd(data *service.PersonFullName) error {
	return validation.ValidateStruct(data,
		validation.Field(&data.Name, validation.Required, validation.Match(regexp.MustCompile(`^[A-Z][a-z]+$`))),
		validation.Field(&data.Surname, validation.Required, validation.Match(regexp.MustCompile(`^[A-Z][a-z]+$`))),
		validation.Field(&data.Patronymic, validation.Match(regexp.MustCompile(`^[A-Z][a-z]+$`))),
	)
}

func (Validator) ValidateDataToGet(data *model.Person) error {
	return validation.ValidateStruct(data,
		validation.Field(&data.Name, validation.Match(regexp.MustCompile(`^[A-Z][a-z]+$`))),
		validation.Field(&data.Surname, validation.Match(regexp.MustCompile(`^[A-Z][a-z]+$`))),
		validation.Field(&data.Patronymic, validation.Match(regexp.MustCompile(`^[A-Z][a-z]+$`))),
		validation.Field(&data.Nationality, validation.Match(regexp.MustCompile(`^[A-Z]{2}$`))),
		validation.Field(&data.Gender, validation.In("female", "male")),
	)
}
func (Validator) ValidateDataToUpdate(data *model.Person) error {
	return validation.ValidateStruct(data,
		validation.Field(&data.Name, validation.Required, validation.Match(regexp.MustCompile(`^[A-Z][a-z]+$`))),
		validation.Field(&data.Surname, validation.Required, validation.Match(regexp.MustCompile(`^[A-Z][a-z]+$`))),
		validation.Field(&data.Patronymic, validation.Match(regexp.MustCompile(`^[A-Z][a-z]+$`))),
		validation.Field(&data.Nationality, validation.Required, validation.Match(regexp.MustCompile(`^[A-Z]{2}$`))),
		validation.Field(&data.Age, validation.Required),
		validation.Field(&data.Gender, validation.Required, validation.In("female", "male")),
	)
}
