package repository

import (
	"database/sql"
)

type SystemRepository interface {
	DBCheck() (bool, error)
}

type systemRepository struct {
	db *sql.DB
}

func NewSystemRepository(db *sql.DB) SystemRepository {
	return &systemRepository{
		db: db,
	}
}

func (r *systemRepository) DBCheck() (bool, error) {
	if err := r.db.Ping(); err == nil {
		return true, nil
	}

	return false, nil
}
