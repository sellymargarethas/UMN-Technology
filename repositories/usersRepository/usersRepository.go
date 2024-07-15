package usersrepository

import (
	"database/sql"
	"fmt"
	"log"
	"umn-technology/config"
	"umn-technology/constants"
	"umn-technology/models"
	"umn-technology/repositories"
	"umn-technology/utils"
	. "umn-technology/utils"
)

type usersRepository struct {
	RepoDB repositories.Repository
}

func NewUsersRepository(repoDB repositories.Repository) usersRepository {
	return usersRepository{
		RepoDB: repoDB,
	}
}

const defineColumnUser = `
	nama, email, username,
	password, createdAt `

func (ctx usersRepository) AddUser(user models.RequestAddUser) (int64, error) {
	var id int64
	query := `
		INSERT INTO users (` + defineColumnUser + `)
		VALUES (
			?,?,?,
			?,NOW())
		RETURNING id`

	query = utils.ReplaceSQL(query, "?")

	err := ctx.RepoDB.DB.QueryRow(query, user.Nama, user.Email, user.Username, user.Password).Scan(&id)

	if err != nil {
		log.Println("Error querying AddUser", err)
		return id, err
	}

	return id, nil
}

func (ctx usersRepository) EditUser(user models.RequestUpdateUser) error {
	query := `
		UPDATE users
		SET 
			nama = ?, email = ?,
			username = ?, updatedAt = NOW()
		WHERE id = ? `

	query = utils.ReplaceSQL(query, "?")

	_, err := ctx.RepoDB.DB.Exec(query, user.Nama, user.Email, user.Username, user.Id)

	if err != nil {
		log.Println("Error querying EditUser", err)
		return err
	}

	return nil
}

func (ctx usersRepository) CheckPassword(request models.RequestUpdatePassword) (bool, error) {
	var storedPassword string
	err := ctx.RepoDB.DB.QueryRow("SELECT password FROM users WHERE id=$1", request.Id).Scan(&storedPassword)
	if err != nil {
		return constants.FALSE_VALUE, err
	}

	decryptedPassword, err := Decrypt(config.E_KEY, storedPassword)
	if err != nil {
		return constants.FALSE_VALUE, err
	}

	if request.OldPassword != decryptedPassword {
		return constants.FALSE_VALUE, fmt.Errorf("password does not match")
	}

	return constants.TRUE_VALUE, nil
}

func (ctx usersRepository) ChangePassword(request models.RequestUpdatePassword) (int64, error) {
	var id int64
	// Update with the new password
	query := `
		UPDATE users SET 
			password=$1, 
			updatedAt=NOW() 
		WHERE id=$2
		RETURNING id`

	err := ctx.RepoDB.DB.QueryRow(query, request.NewPassword, request.Id).Scan(&id)
	if err != nil {
		log.Println("Error querying UpdateUserPassword", err)
		return id, err
	}

	return id, nil
}

func (ctx usersRepository) RemoveUser(id int64) error {
	_, err := ctx.RepoDB.DB.Exec(`
		UPDATE users SET
			deletedAt= NOW() 
		WHERE id=$1`, id)

	if err != nil {
		log.Println("Error querying RemoveUser: ", err)
		return err
	}

	return err
}

func (ctx usersRepository) GetAllUser(request models.RequestList) (users []models.Users, err error) {
	var args []interface{}
	query := `
		SELECT 
			id, nama, email, 
			username, password, createdAt, updatedAt, deletedAt
		FROM users
		WHERE deletedAt IS NULL 
	`

	if request.Keyword != "" {
		query += ` AND (CAST( users.id AS TEXT) ILIKE '%' || ? || '%'
		OR users.nama ILIKE '%' || ? || '%'
		OR users.email ILIKE '%' || ? || '%'
		OR users.username ILIKE '%' || ? || '%')
		`
		args = append(args, request.Keyword, request.Keyword, request.Keyword, request.Keyword)
	}

	orderby := fmt.Sprintf(request.OrderBy)
	order := fmt.Sprintf(request.Order)

	query += ` ORDER BY ` + orderby + ` ` + order

	query = ReplaceSQL(query, "?")

	rows, err := ctx.RepoDB.DB.Query(query, args...)

	if err != nil {
		log.Println("Error querying GetUser", err)
	}

	data, err := usersDto(rows)

	if err != nil {
		return data, err
	}

	return data, nil
}

func (ctx usersRepository) GetUsersListByIndex(user models.Users) ([]models.Users, error) {
	var result []models.Users
	var args []interface{}
	query := `
		SELECT id, ` + defineColumnUser + `, updatedAt, deletedAt
		FROM users
		WHERE deletedAt IS NULL `

	if user.Nama != constants.EMPTY_VALUE {
		query += `AND nama = ? `
		args = append(args, user.Nama)
	}

	if user.Email != constants.EMPTY_VALUE {
		query += `AND email = ? `
		args = append(args, user.Email)
	}

	if user.Username != constants.EMPTY_VALUE {
		query += `AND username = ? `
		args = append(args, user.Username)
	}

	query = utils.ReplaceSQL(query, "?")
	rows, err := ctx.RepoDB.DB.Query(query, args...)
	if err != nil {
		return result, err
	}

	defer rows.Close()
	data, err := usersDto(rows)
	if err != nil {
		return result, err
	}

	return data, nil
}

func (ctx usersRepository) IsUserExistsByIndex(user models.Users) (models.Users, bool, error) {
	var result models.Users
	var args []interface{}

	query := `
		SELECT id, nama, email, username,
		       password, createdAt, updatedAt, deletedAt
		FROM users
		WHERE deletedAt IS NULL `

	if user.Email != constants.EMPTY_VALUE {
		query += ` AND (email = $1`
		args = append(args, user.Email)
	}

	if user.Username != constants.EMPTY_VALUE {
		query += ` OR username = $2 `
		args = append(args, user.Username)
	}

	query += `)`

	rows, err := ctx.RepoDB.DB.Query(query, args...)

	if err != nil {
		return result, constants.FALSE_VALUE, err
	}
	defer rows.Close() // Ensure rows are closed after processing

	data, err := usersDto(rows)
	if err != nil {
		return result, constants.FALSE_VALUE, err
	}

	if len(data) == 0 {
		return result, constants.FALSE_VALUE, nil
	}

	return data[0], constants.TRUE_VALUE, nil
}

func (ctx usersRepository) IsUserExistsByID(id int64) (models.Users, bool, error) {
	var result models.Users
	query := `
		SELECT id, nama, email, username,
		       password, createdAt, updatedAt, deletedAt
		FROM users
		WHERE deletedAt IS NULL AND id=$1`

	rows, err := ctx.RepoDB.DB.Query(query, id)
	if err != nil {
		return result, constants.FALSE_VALUE, err
	}
	defer rows.Close() // Ensure rows are closed after processing

	data, err := usersDto(rows)
	if err != nil {
		return result, constants.FALSE_VALUE, err
	}

	if len(data) == 0 {
		return result, constants.FALSE_VALUE, nil
	}

	return data[0], constants.TRUE_VALUE, nil
}

func usersDto(rows *sql.Rows) ([]models.Users, error) {
	var result []models.Users

	for rows.Next() {
		var val models.Users
		err := rows.Scan(&val.Id, &val.Nama, &val.Email, &val.Username, &val.Password, &val.CreatedAt, &val.UpdatedAt, &val.DeletedAt)
		if err != nil {
			return result, err
		}
		result = append(result, val)
	}

	if err := rows.Err(); err != nil {
		return result, err
	}

	return result, nil
}
