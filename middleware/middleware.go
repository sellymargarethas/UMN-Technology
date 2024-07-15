package middleware

import (
	"net/http"
	"strings"
	"time"
	"umn-technology/utils"

	"github.com/labstack/echo"
)

func AuthenticateMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// Extract access token from Authorization header
		authHeader := ctx.Request().Header.Get("Authorization")
		if authHeader == "" {
			return ctx.JSON(http.StatusUnauthorized, "Missing Authorization Token")
		}

		// Check if the authorization header starts with "Bearer "
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			return ctx.JSON(http.StatusUnauthorized, "Invalid Authorization Token: Missing Bearer prefix")
		}

		// Extract the access token without the "Bearer " prefix
		accessToken := authHeader[len(bearerPrefix):]

		// Validate the access token
		claims, err := utils.ValidateToken(accessToken)
		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, "Invalid Authorization Token: "+err.Error())
		}

		// Check if access token is expired
		expTime := time.Unix(int64(claims["exp"].(float64)), 0)
		if expTime.Before(time.Now()) {
			// Access token is expired, try to refresh it using the refresh token
			refreshToken := ctx.Request().Header.Get("Refresh-Token")
			if refreshToken == "" {
				return ctx.JSON(http.StatusUnauthorized, "Missing Refresh-Token Header")
			}

			newAccessToken, newRefreshToken, err := utils.RefreshAccessToken(refreshToken)
			if err != nil {
				return ctx.JSON(http.StatusUnauthorized, "Failed to refresh access token: "+err.Error())
			}

			// Set the new access token in response headers or cookies
			ctx.Response().Header().Set("Authorization", "Bearer "+newAccessToken)
			ctx.Response().Header().Set("RefreshToken", newRefreshToken)
		}

		// Proceed to the next handler
		return next(ctx)
	}
}

func RefreshTokenHandler(ctx echo.Context) error {
	// Extract refresh token from request body or headers
	refreshToken := ctx.FormValue("refreshToken")
	if refreshToken == "" {
		return ctx.JSON(http.StatusBadRequest, "Missing refreshToken parameter")
	}

	// Attempt to refresh the access token
	newAccessToken, newRefreshToken, err := utils.RefreshAccessToken(refreshToken)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, "Failed to refresh access token")
	}

	// Set the new access token in the response headers or cookies
	ctx.Response().Header().Set("Authorization", newAccessToken)
	ctx.Response().Header().Set("RefreshToken", newRefreshToken)

	// Return success response
	return ctx.JSON(http.StatusOK, "Access token refreshed successfully")
}
