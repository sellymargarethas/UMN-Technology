package main

import (
	"umn-technology/app"
	"umn-technology/config"
	"umn-technology/constants"

	"context"
	"errors"
	"fmt"
	. "umn-technology/utils"

	"net/http"
	"strconv"
	"umn-technology/repositories"
	"umn-technology/routes"

	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/go-playground/validator.v9"
	id_translations "gopkg.in/go-playground/validator.v9/translations/id"
)

// CustomValidator adalah
type CustomValidator struct {
	validator  *validator.Validate
	translator ut.Translator
}

// Passing Variable
var (
	uni         *ut.UniversalTranslator
	echoHandler echo.Echo
)

var ctx = context.Background()

// Custom Validator and translation
func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.validator.Struct(i)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, row := range errs {
			return errors.New(row.Translate(cv.translator))
		}
	}

	return cv.validator.Struct(i)
}

func main() {
	if err := config.OpenConnection(); err != nil {
		panic(fmt.Sprintf("Open Connection Faild: %s", err.Error()))
	}
	defer config.CloseConnectionDB()

	// Connection database
	DB := config.DBConnection()

	// Configuration Repository
	repo := repositories.NewRepository(DB, ctx)

	// Configuration Repository and Services
	services := app.SetupApp(DB, repo)

	// Routing API
	routes.RoutesApi(echoHandler, services)
	echoHandler.Use(middleware.Logger())

	echoHandler.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: constants.TRUE_VALUE,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	port := fmt.Sprintf(":%s", config.GetEnv("APP_PORT"))
	echoHandler.Logger.Fatal(echoHandler.Start(port))
}

func init() {
	e := echo.New()
	echoHandler = *e
	validateCustom := validator.New()

	id := id.New()
	uni = ut.New(id, id)
	trans, _ := uni.GetTranslator("id")
	id_translations.RegisterDefaultTranslations(validateCustom, trans)
	e.Validator = &CustomValidator{validator: validateCustom, translator: trans}

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		report, ok := err.(*echo.HTTPError)
		if !ok {
			report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		result := ResponseJSON(constants.FALSE_VALUE, strconv.Itoa(report.Code), err.Error(), nil)
		c.Logger().Error(report)
		c.JSON(report.Code, result)
	}
}
