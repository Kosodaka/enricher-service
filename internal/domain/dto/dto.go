package dto

type AddPersonDTO struct {
	Id         int    `json:"id,string"`
	Name       string `json:"name" db:"name"`
	Surname    string `json:"surname" db:"surname"`
	Patronymic string `json:"patronymic" db:"patronymic"`
}
