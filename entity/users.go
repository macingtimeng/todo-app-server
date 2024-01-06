package entity

import (
	"strings"
	"todo-app/infra/config"
	"todo-app/pkg/errs"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Email    string `gorm:"unique"`
	Password string
	Todos    []Todo
}

func (u *User) parseToken(tokenString string) (*jwt.Token, errs.Error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.NewUnauthenticatedError("invalid token")
		}

		return []byte(config.AppConfig().JwtSecretKey), nil
	})

	if err != nil {
		return nil, errs.NewUnauthenticatedError("invalid token")
	}

	return token, nil
}

func (u *User) bindTokenToEntity(mapClaims jwt.MapClaims) errs.Error {

	if id, ok := mapClaims["id"].(float64); !ok {
		return errs.NewUnauthenticatedError("invalid token")
	} else {
		u.ID = uint(id)
	}

	if email, ok := mapClaims["email"].(string); !ok {
		return errs.NewUnauthenticatedError("invalid token")
	} else {
		u.Email = email
	}

	return nil
}

func (u *User) ValidateToken(bearerToken string) errs.Error {

	isBearer := strings.HasPrefix(bearerToken, "Bearer")

	if !isBearer {
		return errs.NewUnauthenticatedError("invalid token")
	}

	splitToken := strings.Split(bearerToken, " ")

	if len(splitToken) != 2 {
		return errs.NewUnauthenticatedError("invalid token")
	}

	tokenString := splitToken[1]
	token, err := u.parseToken(tokenString)

	if err != nil {
		return err
	}

	mapClaims := jwt.MapClaims{}

	claims, ok := token.Claims.(jwt.MapClaims)

	// validate token
	if !ok || !token.Valid {
		return errs.NewUnauthenticatedError("invalid token")
	}

	mapClaims = claims

	return u.bindTokenToEntity(mapClaims)
}

func (u *User) claim() jwt.MapClaims {
	return jwt.MapClaims{
		"id":    u.ID,
		"email": u.Email,
	}
}

func (u *User) signToken(claim jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, _ := token.SignedString([]byte(config.AppConfig().JwtSecretKey))
	return tokenString
}

func (u *User) GenerateToken() string {
	return u.signToken(u.claim())
}

func (u *User) HashPassword() error {

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	u.Password = string(hashPassword)

	return nil
}

func (u *User) CompareHashPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
