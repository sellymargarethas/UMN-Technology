package repositories

import (
	"umn-technology/models"
)

type UsersRepository interface {
	AddUser(user models.RequestAddUser) (int64, error)
	EditUser(user models.RequestUpdateUser) error
	CheckPassword(request models.RequestUpdatePassword) (bool, error)
	ChangePassword(request models.RequestUpdatePassword) (int64, error)
	RemoveUser(id int64) error
	GetAllUser(request models.RequestList) (users []models.Users, err error)
	GetUsersListByIndex(user models.Users) ([]models.Users, error)
	IsUserExistsByIndex(user models.Users) (models.Users, bool, error)
	IsUserExistsByID(id int64) (models.Users, bool, error)
}

type LoginRepository interface {
	CheckLogin(data models.Login) (hash string, err error)
	LoginReturn(data models.Login) (user models.LoginResponse, err error)
}
