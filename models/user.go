package models

type Users struct {
	Id        int64   `json:"id"`
	Nama      string  `json:"nama"`
	Email     string  `json:"email"`
	Username  string  `json:"username"`
	Password  string  `json:"password"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt *string `json:"updatedAt"`
	DeletedAt *string `json:"deletedAt"`
}

type RequestAddUser struct {
	Nama     string `json:"nama" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RequestUpdateUser struct {
	Id       int64  `json:"id" validate:"required"`
	Nama     string `json:"nama" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Username string `json:"username" validate:"required"`
}

type RequestUpdatePassword struct {
	Id          int64  `json:"id" validate:"required"`
	OldPassword string `json:"oldPassword" validate:"required,min=8"`
	NewPassword string `json:"newPassword" validate:"required,min=8"`
}

type RequestDeleteUser struct {
	Id int64 `json:"id" validate:"required"`
}
