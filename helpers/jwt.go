package helpers

import (
	"errors"
	"time"

	"github.com/BlenDMinh/dutgrad-server/configs"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTToken(userID uint) (string, time.Time, error) {
	config := configs.GetEnv()
	jwtSecret := config.JwtSecret
	if jwtSecret == "" {
		jwtSecret = "your_default_jwt_secret"
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

func VerifyJWTToken(tokenString string) (uint, error) {
	config := configs.GetEnv()
	jwtSecret := config.JwtSecret
	if jwtSecret == "" {
		jwtSecret = "your_default_jwt_secret"
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid user ID in token")
	}

	return uint(userIDFloat), nil
}
