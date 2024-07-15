package routes

import (
	"net/http"
	"umn-technology/config"
	"umn-technology/constants"
	"umn-technology/handlers"
	umnMiddleware "umn-technology/middleware"
	"umn-technology/services"

	"github.com/labstack/echo/middleware"

	usersService "umn-technology/services/usersService"

	"github.com/labstack/echo"
)

// Routes API
func RoutesApi(e echo.Echo, usecaseSvc services.UsecaseService) {
	public := e.Group("")

	private := e.Group("")
	private.Use(middleware.JWT([]byte(config.GetEnv("JWT_KEY"))))
	private.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: constants.TRUE_VALUE,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	loginGroup := public.Group("/login")
	loginSvc := usersService.NewloginService(usecaseSvc)
	loginGroup.POST("/login", loginSvc.Login)
	loginGroup.POST("/refreshToken", handlers.RefreshTokenHandler)

	userGroupPublic := public.Group("/register")
	userPublicSvc := usersService.NewUsersService(usecaseSvc)
	userGroupPublic.POST("/register", userPublicSvc.InsertUser)

	userGroup := private.Group("/user")
	userGroup.Use(umnMiddleware.AuthenticateMiddleware)
	userSvc := usersService.NewUsersService(usecaseSvc)
	userGroup.POST("/list", userSvc.GetAllUser)
	userGroup.POST("/edit", userSvc.EditUser)
	userGroup.POST("/remove", userSvc.DeleteUser)
	userGroup.POST("/get", userSvc.GetUser)
	userGroup.POST("/updatePassword", userSvc.EditUserPassword)
}
