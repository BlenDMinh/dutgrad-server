package helpers

import (
	"errors"
	"time"

	"github.com/BlenDMinh/dutgrad-server/configs"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTToken(userID uint) (string, time.Time, error) {
	// Get JWT secret from environment variable or use a default one
	config := configs.GetEnv()
	jwtSecret := config.JwtSecret
	if jwtSecret == "" {
		jwtSecret = "your_default_jwt_secret" // In production, never use a hardcoded secret
	}

	// Set token expiration time
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
		"iat":     time.Now().Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

// VerifyToken validates a JWT token and returns the associated user ID
func VerifyJWTToken(tokenString string) (uint, error) {
	// Get JWT secret from environment variable or use a default one
	config := configs.GetEnv()
	jwtSecret := config.JwtSecret
	if jwtSecret == "" {
		jwtSecret = "your_default_jwt_secret" // In production, never use a hardcoded secret
	}

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return 0, err
	}

	// Verify token is valid
	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	// Extract user ID from claims
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
