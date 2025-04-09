package utils

import (
	"errors"
	"time"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/golang-jwt/jwt/v5"
)

// secret key to sign jwt, in production use environment variable
var jwtKey = []byte("hotelqu_secret_key") //hotel_qu_secret_key

type JWTClaims struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken creates JWT tokens for employees
func GenerateToken(employee models.Employee) (string, error) {
	// Set token expiration time (24 hours)
	expirationTime := time.Now().Add(24 * time.Hour)
	
	claims := &JWTClaims{
		Id:    employee.Id,
		Email: employee.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

func ValidateToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	
	return claims, nil
}