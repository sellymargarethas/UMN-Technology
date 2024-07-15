package usersservice

import (
	"log"
	"net/http"
	"umn-technology/config"
	"umn-technology/constants"
	"umn-technology/models"
	"umn-technology/services"
	. "umn-technology/utils"

	"github.com/labstack/echo"
)

type loginService struct {
	Service services.UsecaseService
}

func NewloginService(service services.UsecaseService) loginService {
	return loginService{
		Service: service,
	}
}

func (svc loginService) Login(ctx echo.Context) error {
	var result models.Response
	var request models.Login

	// Validate and bind the login request
	if err := BindValidateStruct(ctx, &request); err != nil {
		log.Println("Error Validate Login: Login", err.Error())
		result = ResponseJSON(constants.FALSE_VALUE, constants.VALIDATE_ERROR_CODE, err.Error(), nil)
		return ctx.JSON(http.StatusBadRequest, result)
	}

	reqLogin := models.Login{
		Username: request.Username,
		Password: request.Password,
	}

	// Check login credentials
	passUser, err := svc.Service.LoginRepo.CheckLogin(reqLogin)
	if err != nil {
		log.Println("Error Login: Check login", err.Error())
		result = ResponseJSON(constants.FALSE_VALUE, constants.VALIDATE_ERROR_CODE, "User Not Exist", nil)
		return ctx.JSON(http.StatusBadRequest, result)
	}

	// Handle user not found
	if passUser == constants.EMPTY_VALUE {
		log.Println("Error Login: User Not Found", err)
		result = ResponseJSON(constants.FALSE_VALUE, constants.VALIDATE_ERROR_CODE, "Wrong Username/Password", nil)
		return ctx.JSON(http.StatusBadRequest, result)
	}

	// Decrypt stored password
	password, err := Decrypt(config.E_KEY, passUser)
	if err != nil {
		log.Println("Error Login-Decrypt : ", err.Error())
		result = ResponseJSON(constants.FALSE_VALUE, constants.FAILED_CODE, "Failed Decrypt Password", err.Error())
		return ctx.JSON(http.StatusBadRequest, result)
	}

	// Check if passwords match
	if password == request.Password {
		// Generate new access token
		token, err := GenerateToken(reqLogin)
		if err != nil {
			log.Println("Error Generating Token: Login", err.Error())
			result = ResponseJSON(constants.FALSE_VALUE, constants.VALIDATE_ERROR_CODE, "Error Generating Token", nil)
			return ctx.JSON(http.StatusBadRequest, result)
		}

		// Generate refresh token
		refreshToken, err := GenerateRefreshToken(reqLogin)
		if err != nil {
			log.Println("Error Generating Refresh Token: Login", err.Error())
			result = ResponseJSON(constants.FALSE_VALUE, constants.VALIDATE_ERROR_CODE, "Error Generating Refresh Token", nil)
			return ctx.JSON(http.StatusBadRequest, result)
		}

		// Retrieve additional login information
		resLogin, err := svc.Service.LoginRepo.LoginReturn(reqLogin)
		if err != nil {
			log.Println("Error Get Login Return: Login", err.Error())
			result = ResponseJSON(constants.FALSE_VALUE, constants.VALIDATE_ERROR_CODE, "Error Getting Information", nil)
			return ctx.JSON(http.StatusBadRequest, result)
		}

		// Include refresh token in the response
		resLogin.Token = token
		resLogin.RefreshToken = refreshToken

		// Construct success response
		result = ResponseJSON(constants.TRUE_VALUE, constants.SUCCESS_CODE, "Berhasil Login", resLogin)

		return ctx.JSON(http.StatusOK, result)
	}

	// Handle incorrect password
	result = ResponseJSON(constants.FALSE_VALUE, constants.VALIDATE_ERROR_CODE, "Wrong Username/Password", nil)
	return ctx.JSON(http.StatusBadRequest, result)
}
