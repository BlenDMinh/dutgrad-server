package helpers

import (
	"errors"
	"time"

	"github.com/BlenDMinh/dutgrad-server/configs"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTToken(userID uint) (string, time.Time, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	token, exp, err := GenerateTokenForPayload(map[string]interface {}{"user_id": userID}, &expirationTime)

	if err != nil {
		return "", time.Time{}, err
	}

	return token, exp, nil
}

func VerifyJWTToken(tokenString string) (uint, error) {
	claims, err := VerifyTokenForPayload(tokenString)
	if err != nil {
		return 0, err
	}

	if claims == nil {
		return 0, errors.New("invalid token claims")
	}

	userIDFloat, ok := (*claims)["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid user ID in token")
	}

	return uint(userIDFloat), nil
}

func GenerateTokenForPayload(payload map[string]interface {}, exp *time.Time) (string, time.Time, error) {
	config := configs.GetEnv()
	jwtSecret := config.JwtSecret
	if jwtSecret == "" {
		jwtSecret = "your_default_jwt_secret"
	}

	var claims jwt.MapClaims

	if exp == nil {
		claims = jwt.MapClaims{
			"payload": payload,
			"iat":     time.Now().Unix(),
		}
	} else {
		claims = jwt.MapClaims{
			"payload": payload,
			"exp":     exp.Unix(),
			"iat":     time.Now().Unix(),
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, *exp, nil	
}

func VerifyTokenForPayload(tokenString string) (*map[string]interface {}, error) {
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
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	payload := claims["payload"].(map[string]interface {})

	return &payload, nil
}
