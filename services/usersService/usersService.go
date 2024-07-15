package usersservice

import (
	"fmt"
	"log"
	"net/http"
	config "umn-technology/config"
	"umn-technology/constants"
	"umn-technology/models"
	"umn-technology/services"
	. "umn-technology/utils"

	"github.com/labstack/echo"
)

type usersService struct {
	Service services.UsecaseService
}

func NewUsersService(service services.UsecaseService) usersService {
	return usersService{
		Service: service,
	}
}

func (svc usersService) InsertUser(ctx echo.Context) error {
	var result models.Response
	var request models.RequestAddUser

	if err := BindValidateStruct(ctx, &request); err != nil {
		log.Println("Error Validate Data:  InsertUser", err.Error())
		result = ResponseJSON(constants.FALSE_VALUE, constants.VALIDATE_ERROR_CODE, "Failed Validate Data", err.Error())
		return ctx.JSON(http.StatusBadRequest, result)
	}
	users := models.Users{
		Username: request.Username,
		Email:    request.Email,
	}

	_, exists, err := svc.Service.UserRepo.IsUserExistsByIndex(users)
	fmt.Println(exists)
	if exists {
		log.Println("Error EditUser - IsUserExistsByIndex : User Already Exists ")
		result = ResponseJSON(constants.FALSE_VALUE, constants.FAILED_CODE, "User Already Exists", nil)
		return ctx.JSON(http.StatusOK, result)
	}

	if err != nil {
		log.Println("Error EditUser- IsUserExistsByIndex : ", err.Error())
		result = ResponseJSON(constants.FALSE_VALUE, constants.FAILED_CODE, "Failed Edit User", nil)
		return ctx.JSON(http.StatusBadRequest, result)
	}

	hash, err := Encrypt(config.E_KEY, request.Password)
	if err != nil {
		log.Println("Error InsertUser-Encrypt : ", err.Error())
		result = ResponseJSON(constants.FALSE_VALUE, constants.FAILED_CODE, "Failed Encrypt Password", err.Error())

		return ctx.JSON(http.StatusBadRequest, result)
	}

	request.Password = string(hash)

	resInsertUser, err := svc.Service.UserRepo.AddUser(request)
	if err != nil {
		log.Println("Error InsertUser-AddUser : ", err.Error())
		result = ResponseJSON(constants.FALSE_VALUE, constants.FAILED_CODE, "Failed Insert User", err.Error())

		return ctx.JSON(http.StatusBadRequest, result)
	}

	result = ResponseJSON(constants.TRUE_VALUE, constants.SUCCESS_CODE, "Success Insert User", resInsertUser)

	return ctx.JSON(http.StatusOK, result)

}

func (svc usersService) GetAllUser(ctx echo.Context) (err error) {
	var result models.ResponseList
	var request models.RequestList
	var orderby string

	if err := BindValidateStruct(ctx, &request); err != nil {
		log.Println("Error Validate Data GetAllUser: ", err.Error())
		result = ResponseListJSON(constants.FALSE_VALUE, constants.VALIDATE_ERROR_CODE, err.Error(), constants.EMPTY_VALUE_INT, nil)
		return ctx.JSON(http.StatusBadRequest, result)
	}

	if request.OrderBy == "id" {
		orderby = "users.id"
	}

	if request.OrderBy == "nama" {
		orderby = "users.nama"
	}

	if request.OrderBy == "email" {
		orderby = "users.email"
	}

	if request.OrderBy == "username" {
		orderby = "users.username"
	}

	if request.OrderBy == constants.EMPTY_VALUE {
		orderby = "users.id"
	}

	request.OrderBy = orderby

	resuser, err := svc.Service.UserRepo.GetAllUser(request)

	if err != nil {
		log.Println("Error GetAllUser- GetAllUser : ", err.Error())
		result = ResponseListJSON(constants.FALSE_VALUE, constants.FAILED_CODE, "Failed Get User", constants.EMPTY_VALUE_INT, nil)
		return ctx.JSON(http.StatusBadRequest, result)
	}
	limitbawah := (request.Page - 1) * request.Limit
	limitatas := request.Limit * request.Page

	if len(resuser) < limitatas {
		limitatas = len(resuser)
	}
	if limitbawah > len(resuser) {
		log.Println("Error GetAllUser-GetAllUser: limit melebihi jumlah data", err)
		result = ResponseListJSON(constants.FALSE_VALUE, constants.FAILED_CODE, "Limit melebihi jumlah data", len(resuser), err)
		return ctx.JSON(http.StatusBadRequest, result)
	}

	userslice := resuser[limitbawah:limitatas]
	if len(userslice) == constants.EMPTY_VALUE_INT {
		log.Println("Error GetAllUser-GetAllUser: limit melebihi jumlah data", err)
		result = ResponseListJSON(constants.FALSE_VALUE, constants.DATA_NOT_FOUND_CODE, "Data Tidak Ditemukan", len(resuser), err)
		return ctx.JSON(http.StatusBadRequest, result)
	}

	log.Println("Reponse GetAllUser-GetAllUser")
	result = ResponseListJSON(constants.TRUE_VALUE, constants.SUCCESS_CODE, "Success Get User", len(resuser), userslice)

	return ctx.JSON(http.StatusOK, result)
}

func (svc usersService) GetUser(ctx echo.Context) (err error) {
	var result models.Response
	var request models.Users

	if err := BindValidateStruct(ctx, &request); err != nil {
		log.Println("Error Validate Data", err.Error())
		result = ResponseJSON(constants.FALSE_VALUE, constants.VALIDATE_ERROR_CODE, "Failed Validate Data", err.Error())
		return ctx.JSON(http.StatusBadRequest, result)
	}

	resGetUser, err := svc.Service.UserRepo.GetUsersListByIndex(request)
	if err != nil {

		log.Println("Error GetUser-GetUsersListByIndex : ", err.Error())
		result = ResponseJSON(constants.FALSE_VALUE, constants.FAILED_CODE, "Failed Get User", nil)

		return ctx.JSON(http.StatusBadRequest, result)
	}

	if len(resGetUser) == 0 {
		log.Println("Error GetUser-GetUsersListByIndex : User Not Found ", err)
		result = ResponseJSON(constants.FALSE_VALUE, constants.DATA_NOT_FOUND_CODE, "User Not Found", nil)
		return ctx.JSON(http.StatusBadRequest, result)
	}
	log.Println("Reponse GetUser-GetUsersListByIndex")
	result = ResponseJSON(constants.TRUE_VALUE, constants.SUCCESS_CODE, "Success Get User", resGetUser)

	return ctx.JSON(http.StatusOK, result)

}

func (svc usersService) EditUser(ctx echo.Context) error {
	var result models.Response
	var request models.RequestUpdateUser
	if err := BindValidateStruct(ctx, &request); err != nil {
		log.Println("Error Validate Data", err.Error())
		result = ResponseJSON(constants.FALSE_VALUE, constants.VALIDATE_ERROR_CODE, "Failed Validate Data", err.Error())
		return ctx.JSON(http.StatusBadRequest, result)
	}

	_, exists, err := svc.Service.UserRepo.IsUserExistsByID(request.Id)
	fmt.Println(exists)
	if !exists {
		log.Println("Error EditUser - IsUserExistsByID : User Not Exists ")
		result = ResponseJSON(constants.FALSE_VALUE, constants.FAILED_CODE, "User Not Exists", nil)
		return ctx.JSON(http.StatusOK, result)
	}

	if err != nil {
		log.Println("Error EditUser- IsUserExistsByID : ", err.Error())
		result = ResponseJSON(constants.FALSE_VALUE, constants.FAILED_CODE, "Failed Edit User", nil)
		return ctx.JSON(http.StatusBadRequest, result)
	}

	err = svc.Service.UserRepo.EditUser(request)
	if err != nil {
		log.Println("Error EditUser- EditUser : ", err.Error())
		result = ResponseJSON(constants.FALSE_VALUE, constants.FAILED_CODE, "Failed Edit User", nil)
		return ctx.JSON(http.StatusBadRequest, result)
	}

	result = ResponseJSON(constants.TRUE_VALUE, constants.SUCCESS_CODE, "Success Edit User", nil)

	return ctx.JSON(http.StatusOK, result)

}
func (svc usersService) EditUserPassword(ctx echo.Context) error {
	var result models.Response
	var request models.RequestUpdatePassword
	if err := BindValidateStruct(ctx, &request); err != nil {
		log.Println("Error Validate Data", err.Error())
		result = ResponseJSON(constants.FALSE_VALUE, constants.VALIDATE_ERROR_CODE, "Failed Validate Data", err.Error())
		return ctx.JSON(http.StatusBadRequest, result)
	}

	_, exists, err := svc.Service.UserRepo.IsUserExistsByID(request.Id)
	if !exists {
		log.Println("Error EditUserPassword - IsUserExistsByID : User Not Exists")
		result = ResponseJSON(constants.FALSE_VALUE, constants.FAILED_CODE, "User Not Exists", nil)
		return ctx.JSON(http.StatusOK, result)
	}
	if err != nil {
		log.Println("Error EditUserPassword- IsUserExistsByID : ", err.Error())
		result = ResponseJSON(constants.FALSE_VALUE, constants.FAILED_CODE, "Failed Edit User", nil)
		return ctx.JSON(http.StatusBadRequest, result)
	}

	passwordMatch, err := svc.Service.UserRepo.CheckPassword(request)
	if !passwordMatch {
		log.Println("Error EditUserPassword - CheckPassword : Old password does not match")
		result = ResponseJSON(constants.FALSE_VALUE, constants.FAILED_CODE, "Old password does not match", nil)
		return ctx.JSON(http.StatusOK, result)
	}
	if err != nil {
		log.Println("Error EditUserPassword- CheckPassword : ", err.Error())
		result = ResponseJSON(constants.FALSE_VALUE, constants.FAILED_CODE, "Failed to check password", nil)
		return ctx.JSON(http.StatusBadRequest, result)
	}

	hash, err := Encrypt(config.E_KEY, request.NewPassword)
	if err != nil {
		log.Println("Error EditUserPassword-Encrypt : ", err.Error())
		result = ResponseJSON(constants.FALSE_VALUE, constants.FAILED_CODE, "Failed Encrypt Password", err.Error())
		return ctx.JSON(http.StatusBadRequest, result)
	}

	request.NewPassword = string(hash)

	resUpdate, err := svc.Service.UserRepo.ChangePassword(request)
	if err != nil {
		log.Println("Error EditUserPassword-ChangePassword : ", err.Error())
		result = ResponseJSON(constants.FALSE_VALUE, constants.FAILED_CODE, "Failed Edit User", nil)
		return ctx.JSON(http.StatusBadRequest, result)
	}
	result = ResponseJSON(constants.TRUE_VALUE, constants.SUCCESS_CODE, "Success Edit User", resUpdate)

	return ctx.JSON(http.StatusOK, result)
}

func (svc usersService) DeleteUser(ctx echo.Context) (err error) {
	var result models.Response
	var request models.RequestDeleteUser
	if err := BindValidateStruct(ctx, &request); err != nil {

		log.Println("Error Validate Data", err.Error())
		result = ResponseJSON(constants.FALSE_VALUE, constants.VALIDATE_ERROR_CODE, "Failed Validate Data", err.Error())
		return ctx.JSON(http.StatusBadRequest, result)
	}
	users := models.Users{
		Id: request.Id,
	}
	_, exists, err := svc.Service.UserRepo.IsUserExistsByID(users.Id)
	if !exists {
		log.Println("Error DeleteUser - IsUserExistsByID : User Not Exists")
		result = ResponseJSON(constants.FALSE_VALUE, constants.FAILED_CODE, "User Not Exists", nil)
		return ctx.JSON(http.StatusOK, result)
	}

	if err != nil {
		log.Println("Error DeleteUser- IsUserExistsByID : ", err.Error())
		result = ResponseJSON(constants.FALSE_VALUE, constants.FAILED_CODE, "Failed Edit User", nil)
		return ctx.JSON(http.StatusBadRequest, result)
	}

	err = svc.Service.UserRepo.RemoveUser(users.Id)
	if err != nil {
		log.Println("Error DeleteUser-RemoveUser : ", err.Error())
		result = ResponseJSON(constants.FALSE_VALUE, constants.FAILED_CODE, "Failed Delete User", nil)

		return ctx.JSON(http.StatusBadRequest, result)
	}

	result = ResponseJSON(constants.TRUE_VALUE, constants.SUCCESS_CODE, "Success Delete User", nil)
	log.Println("Reponse Delete User- Delete User", result)

	return ctx.JSON(http.StatusOK, result)

}
