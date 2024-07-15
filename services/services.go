package services

import (
	"database/sql"

	"umn-technology/repositories"
)

type UsecaseService struct {
	RepoDB    *sql.DB
	UserRepo  repositories.UsersRepository
	LoginRepo repositories.LoginRepository
}

func NewUsecaseService(
	repoDB *sql.DB,
	UserRepo repositories.UsersRepository,
	LoginRepo repositories.LoginRepository,
) UsecaseService {
	return UsecaseService{
		RepoDB:    repoDB,
		UserRepo:  UserRepo,
		LoginRepo: LoginRepo,
	}
}
