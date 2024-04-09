package repository

import (
	"context"
	"fmt"
	"github.com/Kosodaka/enricher-service/internal/domain/model"
	"github.com/jmoiron/sqlx"
)

type PersonRepository struct {
	db *sqlx.DB
}

func NewPersonPostgres(db *sqlx.DB) *PersonRepository {
	return &PersonRepository{
		db: db,
	}
}
func (r *PersonRepository) AddPerson(ctx context.Context, data *model.Person) (int, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	stmt := `INSERT INTO person (name, surname, patronymic, age, gender, nationality) 
			VALUES (:name, :surname, :patronymic, :age, :gender, :nationality) RETURNING id`

	var id int
	insertStmt, err := tx.PrepareNamedContext(ctx, stmt)
	if err != nil {
		return 0, err
	}

	err = insertStmt.GetContext(ctx, &id, data)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PersonRepository) GetPerson(ctx context.Context, id int) (*model.Person, error) {
	stmt := "SELECT id, name, surname, patronymic, age, gender, nationality FROM person WHERE id = $1"
	person := &model.Person{}
	err := r.db.QueryRowxContext(ctx, stmt, id).StructScan(person)
	if err != nil {
		return nil, fmt.Errorf("%s: no such user", err)
	}

	return person, nil
}

func (r *PersonRepository) GetPersons(ctx context.Context, data *model.Person) ([]model.Person, error) {
	stmt := `SELECT id, name, surname, patronymic, age, gender, nationality FROM person WHERE name = name AND surname = surname 
             AND patronymic = patronymic AND age = age AND gender = gender AND nationality = nationality`
	persons := []model.Person{}
	rows, err := r.db.NamedQueryContext(ctx, stmt, data)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		person := &model.Person{}
		if err := rows.StructScan(person); err != nil {
			return nil, err
		}
		persons = append(persons, *person)
	}

	return persons, nil
}

func (r *PersonRepository) UpdatePerson(ctx context.Context, data *model.Person) error {
	tx, err := r.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := "UPDATE person SET name = :name,surname = :surname,patronymic = :patronymic,age = :age,gender = :gender, nationality = :nationality WHERE id = :id"
	updateStmt, err := tx.PrepareNamedContext(ctx, stmt)
	if err != nil {
		return err
	}

	result, err := updateStmt.Exec(data)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("%s: no such user", err)
	}

	return tx.Commit()
}

func (r *PersonRepository) DeletePerson(ctx context.Context, id int) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := "DELETE FROM person WHERE id = $1"

	deleteStmt, err := tx.PrepareContext(ctx, stmt)
	if err != nil {
		return err
	}

	result, err := deleteStmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: no such user", err)
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}
	if rows == 0 {
		fmt.Errorf("%s: no such user", err)
	}

	return tx.Commit()
}
