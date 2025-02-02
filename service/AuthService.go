package service

import (
	"encoding/base64"
	"fmt"
	"jds-test/config"
	"jds-test/model"
	"jds-test/model/constant"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	config config.Config
}

func RegAuthService(
	config config.Config,
) AuthService {
	return AuthService{
		config: config,
	}
} 

// @Summary Register a new user
// @Description Register a new user with NIK and Role
// @Tags auth
// @Accept json
// @Produce json
// @Param user body model.User true "User data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /jds-test/auth/register [post]
func (c AuthService) RegisterUser(e echo.Context) error {
	var (
		user model.User
		response model.ResponseReturn
	)

	err := e.Bind(&user)
	if err != nil {
		response.Code = http.StatusBadRequest
		response.Message = err.Error()
		return e.JSON(http.StatusBadRequest, response)
	}

	err = c.userValidation(user)
	if err != nil {
		response.Code = http.StatusBadRequest
		response.Message = err.Error()
		return e.JSON(http.StatusBadRequest, response)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Code = http.StatusBadRequest
		response.Message = "Failed to hashed password"
		return e.JSON(http.StatusBadRequest, response)
	}

	encodePassword := base64.StdEncoding.EncodeToString(hashedPassword)

	user.Password = encodePassword
	response.Code = http.StatusOK
	response.Message = "ok"
	response.Data = user

	// secret := generateJwtSecret("jds test")
	// fmt.Println(secret)

	return e.JSON(http.StatusOK, response)
}

// @Summary Login a user
// @Description Login a user with NIK and Password
// @Tags auth
// @Accept json
// @Produce json
// @Param login body model.User true "Login data"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
func (c AuthService) Login(e echo.Context) error {
	var (
		user model.User
		response model.ResponseReturn
	)
	

	err := e.Bind(&user)
	if err != nil {
		response.Code = http.StatusBadRequest
		response.Message = err.Error()
		return e.JSON(http.StatusBadRequest, response)
	}

	err = c.userValidation(user)
	if err != nil {
		response.Code = http.StatusBadRequest
		response.Message = err.Error()
		return e.JSON(http.StatusBadRequest, response)
	}

	decodeDummyPassword, err := base64.StdEncoding.DecodeString(constant.PASSWORDHASH_TEST)
	if err != nil {
		response.Code = http.StatusBadRequest
		response.Message = err.Error()
		return e.JSON(http.StatusBadRequest, response)
	}

	err = bcrypt.CompareHashAndPassword(decodeDummyPassword, []byte(constant.PASSWORDTEST))
	if err != nil {
		fmt.Printf("Error CompareHashAndPassword, %s", err.Error())
		response.Code = http.StatusBadRequest
		response.Message = err.Error()
		return e.JSON(http.StatusBadRequest, response)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.TokenClaim{
		Nik: user.Nik,
		Role: user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	})

	jwtSecret, err := base64.StdEncoding.DecodeString(c.config.MainConfig.JwtSecret)
	if err != nil {
		fmt.Printf("Error base64.StdEncoding.DecodeString, %s", err.Error())
		response.Code = http.StatusBadRequest
		response.Message = err.Error()
		return e.JSON(http.StatusBadRequest, response)
	}

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		fmt.Printf("Error signing, %s", err.Error())
		response.Code = http.StatusBadRequest
		response.Message = err.Error()
		return e.JSON(http.StatusBadRequest, response)
	}

	user.Token = tokenString
	response.Code = http.StatusOK
	response.Message = "ok"
	response.Data = user

	return e.JSON(http.StatusOK, response)
}

func (c AuthService) userValidation(user model.User) (err error) {
	if user.Nik == ""  || user.Password == "" || user.Role == "" {
		return fmt.Errorf("Nik, password or role must be provided")
	}

	if len(user.Nik) != 16 {
		return fmt.Errorf("Invalid NIK, NIK must be 16 characters")
	}

	_, err = strconv.Atoi(user.Nik)
	if err != nil {
		return fmt.Errorf("Invalid NIK, NIK should be number")
	}


	if len(user.Password) != 6 {
		return fmt.Errorf("Invalid password, Pasword must be 6 characters")
	}

	return err
}

func generateJwtSecret(value string) (string) {
	bytes := []byte(value)
	secret := base64.StdEncoding.EncodeToString(bytes)
	return secret
}

func (c AuthService) ProtectedEndpoint(e echo.Context) error {
	claims := e.Get("claims").(*model.TokenClaim)
	return e.JSON(http.StatusOK, model.ResponseReturn{Code: http.StatusOK, Message: "Ok", Data: claims})
}