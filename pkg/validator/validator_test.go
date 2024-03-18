package validator

import (
	"github.com/Kosodaka/enricher-service/internal/domain/dto"
	domainErr "github.com/Kosodaka/enricher-service/internal/domain/errors"
	"github.com/Kosodaka/enricher-service/internal/domain/model"
	validation "github.com/go-ozzo/ozzo-validation"
	"reflect"
	"testing"
)

func TestValidateDataToAdd(t *testing.T) {
	validator := NewValidator()

	cases := []struct {
		name   string
		data   *dto.AddPersonDTO
		expErr error
	}{
		{
			name: "valid_data",
			data: &dto.AddPersonDTO{
				Name:       "Dmitriy",
				Surname:    "Ushakov",
				Patronymic: "Vasilevich",
			},
			expErr: nil,
		},
		{
			name: "valid_without_patronymic",
			data: &dto.AddPersonDTO{
				Name:       "Dmitriy",
				Surname:    "Ushakov",
				Patronymic: "",
			},
			expErr: nil,
		},
		{
			name: "invalid_name_lowercase",
			data: &dto.AddPersonDTO{
				Name:    "dmitriy",
				Surname: "Ushakov",
			},
			expErr: validation.Errors{"name": domainErr.InvalidData},
		},
		{
			name: "invalid_name_mixed_uppercase",
			data: &dto.AddPersonDTO{
				Name:    "DmitRiy",
				Surname: "Ushakov",
			},
			expErr: validation.Errors{"name": domainErr.InvalidData},
		},
		{
			name: "invalid_empty_name_field",
			data: &dto.AddPersonDTO{
				Name:    "",
				Surname: "Ushakov",
			},
			expErr: validation.Errors{"name": domainErr.EmptyField},
		},
		{
			name: "invalid_surname_lowercase",
			data: &dto.AddPersonDTO{
				Name:    "Dmitriy",
				Surname: "ushakov",
			},
			expErr: validation.Errors{"surname": domainErr.InvalidData},
		},
		{
			name: "invalid_surname_mixed_uppercase",
			data: &dto.AddPersonDTO{
				Name:    "Dmitriy",
				Surname: "UsHakov",
			},
			expErr: validation.Errors{"surname": domainErr.InvalidData},
		},
		{
			name: "invalid_empty_surname_field",
			data: &dto.AddPersonDTO{
				Name:    "Dmitriy",
				Surname: "",
			},
			expErr: validation.Errors{"surname": domainErr.EmptyField},
		},
	}

	for _, testCases := range cases {
		t.Run(testCases.name, func(t *testing.T) {
			err := validator.ValidateDataToAdd(testCases.data)
			if !reflect.DeepEqual(err, testCases.expErr) {
				t.Errorf("got %v, want %v", err, testCases.expErr)
			}
		})
	}
}
func TestValidateDataToGet(t *testing.T) {
	validator := NewValidator()

	cases := []struct {
		name   string
		data   *model.Person
		expErr error
	}{
		{
			name: "valid_data",
			data: &model.Person{
				Name:       "Dmitriy",
				Surname:    "Ushakov",
				Patronymic: "Vasilevich",
			},
			expErr: nil,
		},
		{
			name: "valid_without_patronymic",
			data: &model.Person{
				Name:       "Dmitriy",
				Surname:    "Ushakov",
				Patronymic: "",
			},
			expErr: nil,
		},
		{
			name: "invalid_name_lowercase",
			data: &model.Person{
				Name:    "dmitriy",
				Surname: "Ushakov",
			},
			expErr: validation.Errors{"name": domainErr.InvalidData},
		},
		{
			name: "invalid_name_mixed_uppercase",
			data: &model.Person{
				Name:    "DmitRiy",
				Surname: "Ushakov",
			},
			expErr: validation.Errors{"name": domainErr.InvalidData},
		},
		{
			name: "invalid_surname_lowercase",
			data: &model.Person{
				Name:    "Dmitriy",
				Surname: "ushakov",
			},
			expErr: validation.Errors{"surname": domainErr.InvalidData},
		},
		{
			name: "invalid_surname_mixed_uppercase",
			data: &model.Person{
				Name:    "Dmitriy",
				Surname: "UsHakov",
			},
			expErr: validation.Errors{"surname": domainErr.InvalidData},
		},
		{
			name: "invalid_nationality_three_chars",
			data: &model.Person{
				Nationality: "USA",
			},
			expErr: validation.Errors{"nationality": domainErr.InvalidData},
		},
		{
			name: "invalid_gender_uppercase",
			data: &model.Person{
				Gender: "Female",
			},
			expErr: validation.Errors{"gender": domainErr.InvalidGender},
		},
	}

	for _, testCases := range cases {
		t.Run(testCases.name, func(t *testing.T) {
			err := validator.ValidateDataToGet(testCases.data)
			if !reflect.DeepEqual(err, testCases.expErr) {
				t.Errorf("got %v, want %v", err, testCases.expErr)
			}
		})
	}
}

func TestValidateDataToUpdate(t *testing.T) {
	validator := NewValidator()

	cases := []struct {
		name   string
		data   *model.Person
		expErr error
	}{
		{
			name: "valid_data",
			data: &model.Person{
				Name:        "Dmitriy",
				Surname:     "Ushakov",
				Patronymic:  "Vasilevich",
				Age:         34,
				Gender:      "male",
				Nationality: "UA",
			},
			expErr: nil,
		},
		{
			name: "valid_without_patronymic",
			data: &model.Person{
				Name:        "Dmitriy",
				Surname:     "Ushakov",
				Patronymic:  "",
				Age:         34,
				Gender:      "male",
				Nationality: "UA",
			},
			expErr: nil,
		},
		{
			name: "invalid_name_lowercase",
			data: &model.Person{
				Name:        "dmitriy",
				Surname:     "Ushakov",
				Age:         34,
				Gender:      "male",
				Nationality: "UA",
			},
			expErr: validation.Errors{"name": domainErr.InvalidData},
		},
		{
			name: "invalid_name_mixed_uppercase",
			data: &model.Person{
				Name:        "DmitRiy",
				Surname:     "Ushakov",
				Age:         34,
				Gender:      "male",
				Nationality: "UA",
			},
			expErr: validation.Errors{"name": domainErr.InvalidData},
		},
		{
			name: "invalid_empty_name_field",
			data: &model.Person{
				Name:        "",
				Surname:     "Ushakov",
				Age:         34,
				Gender:      "male",
				Nationality: "UA",
			},
			expErr: validation.Errors{"name": domainErr.EmptyField},
		},
		{
			name: "invalid_surname_lowercase",
			data: &model.Person{
				Name:        "Dmitriy",
				Surname:     "ushakov",
				Age:         34,
				Gender:      "male",
				Nationality: "UA",
			},
			expErr: validation.Errors{"surname": domainErr.InvalidData},
		},
		{
			name: "invalid_surname_mixed_uppercase",
			data: &model.Person{
				Name:        "Dmitriy",
				Surname:     "UsHakov",
				Age:         34,
				Gender:      "male",
				Nationality: "UA",
			},
			expErr: validation.Errors{"surname": domainErr.InvalidData},
		},
		{
			name: "invalid_empty_surname_field",
			data: &model.Person{
				Name:        "Dmitriy",
				Surname:     "",
				Age:         34,
				Gender:      "male",
				Nationality: "UA",
			},
			expErr: validation.Errors{"surname": domainErr.EmptyField},
		},
		{
			name: "invalid_nationality_three_chars",
			data: &model.Person{
				Name:        "Dmitriy",
				Surname:     "Ushakov",
				Age:         34,
				Gender:      "male",
				Nationality: "USA",
			},
			expErr: validation.Errors{"nationality": domainErr.InvalidData},
		},
		{
			name: "invalid_gender_uppercase",
			data: &model.Person{
				Name:        "Dmitriy",
				Surname:     "Ushakov",
				Age:         34,
				Gender:      "Male",
				Nationality: "UA",
			},
			expErr: validation.Errors{"gender": domainErr.InvalidGender},
		},
	}

	for _, testCases := range cases {
		t.Run(testCases.name, func(t *testing.T) {
			err := validator.ValidateDataToUpdate(testCases.data)
			if !reflect.DeepEqual(err, testCases.expErr) {
				t.Errorf("got %v, want %v", err, testCases.expErr)
			}
		})
	}
}
