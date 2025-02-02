package utils

import (
	"encoding/base64"
	"fmt"
	"jds-test/config"
	"jds-test/model"
	"jds-test/model/constant"
	"net/http"
	"strings"

	_ "jds-test/docs"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Engine struct {
	Echo *echo.Echo
}

// @title Authentication API
// @version 1.0.0
// @description This is a authentication API.
// @BasePath /api/v1
func ReqEngine() Engine {
	engine := echo.New()

	engine.GET("/swagger/*", echoSwagger.WrapHandler)


	return Engine{Echo: engine}
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		response := model.ResponseReturn{}
		jwtSecret, err := base64.StdEncoding.DecodeString(config.Init().MainConfig.JwtSecret)
		if err != nil {
			fmt.Printf("Error base64.StdEncoding.DecodeString, %s", err.Error())
			response.Code = http.StatusBadRequest
			response.Message = err.Error()
			return e.JSON(http.StatusBadRequest, response)
		}

		tokenString := e.Request().Header.Get("Authorization")
		if tokenString == "" {
			response.Code = http.StatusUnauthorized
			response.Message = "Authorization header is required"
			return e.JSON(http.StatusUnauthorized, response)
		}

		fmt.Println(tokenString)

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		token, err := jwt.ParseWithClaims(tokenString, &model.TokenClaim{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		

		if err != nil {
			response.Code = http.StatusUnauthorized
			response.Message = "Parsing token failed"
			return e.JSON(http.StatusUnauthorized, response)
		}

		if !token.Valid {
			response.Code = http.StatusUnauthorized
			response.Message = "Invalid Token"
			return e.JSON(http.StatusUnauthorized, response)
		}

		claims := token.Claims.(*model.TokenClaim)
		e.Set("claims", claims)

		return next(e)
	}
}

func AdminRoleMiddleware() echo.HandlerFunc {
	return func(e echo.Context) error {
		claims := e.Get("claims").(*model.TokenClaim)
		if claims.Role == constant.ROLE_TYPE_ADMIN {
			return e.JSON(http.StatusForbidden, model.ResponseReturn{Code: http.StatusForbidden, Message:"Admin role required"})
		}

		return nil
	}
}