package usecase

import (
	"errors"
	"net/http"
	"ot-recorder/app/response"
	"ot-recorder/app/system/repository"
)

type SystemUsecase interface {
	GetHealth() error
}

type systemUsecase struct {
	repo repository.SystemRepository
}

func NewSystemUsecase(repo repository.SystemRepository) SystemUsecase {
	return &systemUsecase{
		repo: repo,
	}
}

func (u *systemUsecase) GetHealth() error {
	// check db
	dbOnline, err := u.repo.DBCheck()
	if err != nil {
		return err
	}

	if !dbOnline {
		return response.WrapError(errors.New("DB is offline"), http.StatusServiceUnavailable)
	}

	return nil
}
