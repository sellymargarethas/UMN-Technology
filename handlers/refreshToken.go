package handlers

import (
	"net/http"
	"umn-technology/models"
	"umn-technology/utils"

	"github.com/labstack/echo"
)

// RefreshTokenHandler handles the refresh token request
func RefreshTokenHandler(ctx echo.Context) error {
	// Parse the refresh token from the request body
	req := new(models.RefreshTokenRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid request payload",
		})
	}

	refreshToken := req.RefreshToken
	if refreshToken == "" {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error": "Missing refreshToken parameter",
		})
	}

	// Attempt to refresh the access token
	newAccessToken, newRefreshToken, err := utils.RefreshAccessToken(refreshToken)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, echo.Map{
			"error": "Failed to refresh access token",
		})
	}

	// Set the new access token in the response headers or cookies
	ctx.Response().Header().Set("Authorization", newAccessToken)
	ctx.Response().Header().Set("RefreshToken", newRefreshToken)

	// Return success response
	return ctx.JSON(http.StatusOK, echo.Map{
		"message": "Access token refreshed successfully",
	})
}
