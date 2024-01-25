package postgres

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Psql struct {
	DSN string
}

func NewPsql(psqlPath string) *Psql {
	return &Psql{
		DSN: psqlPath,
	}
}

func (p *Psql) GetDb() (*sqlx.DB, error) {
	return sqlx.Open("postgres", p.DSN)
}
