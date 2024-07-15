package usersrepository

import (
	"log"
	"umn-technology/models"
	"umn-technology/repositories"
	. "umn-technology/utils"
)

type loginRepository struct {
	RepoDB repositories.Repository
}

func NewloginRepository(repoDB repositories.Repository) loginRepository {
	return loginRepository{
		RepoDB: repoDB,
	}
}

func (ctx loginRepository) CheckLogin(data models.Login) (hash string, err error) {
	query := `SELECT password FROM users WHERE username=?`

	query = ReplaceSQL(query, "?")
	err = ctx.RepoDB.DB.QueryRow(query, data.Username).Scan(&hash)

	if err != nil {
		log.Println("Error querying: CheckLogin: ", err)
	}
	return hash, err

}

func (ctx loginRepository) LoginReturn(data models.Login) (user models.LoginResponse, err error) {
	query := `
		SELECT users.id, users.nama 
		FROM users
		WHERE users.username=$1`
	err = ctx.RepoDB.DB.QueryRow(query, data.Username).Scan(&user.Id, &user.Nama)
	if err != nil {
		log.Println("Error querying LoginReturn: ", err)
		return user, err
	}
	return
}
