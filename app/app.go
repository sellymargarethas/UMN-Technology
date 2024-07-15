package app

import (
	"database/sql"

	"umn-technology/repositories"

	usersrepository "umn-technology/repositories/usersRepository"

	"umn-technology/services"
)

func SetupApp(DB *sql.DB, repo repositories.Repository) services.UsecaseService {
	usersRepo := usersrepository.NewUsersRepository(repo)
	loginRepository := usersrepository.NewloginRepository(repo)

	usecaseSvc := services.NewUsecaseService(
		DB, usersRepo, loginRepository,
	)

	return usecaseSvc
}
