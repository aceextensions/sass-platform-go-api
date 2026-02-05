package service

import (
	"time"

	"github.com/aceextension/core/config"
	"github.com/aceextension/identity/dto"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateAccessToken(payload dto.TokenPayload) (string, error) {
	claims := jwt.MapClaims{
		"userId":   payload.UserID.String(),
		"role":     payload.Role,
		"exp":      time.Now().Add(time.Hour * 1).Unix(), // 1 hour
		"iat":      time.Now().Unix(),
	}

	if payload.TenantID != nil {
		claims["tenantId"] = payload.TenantID.String()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GlobalConfig.JWTSecret))
}

func GenerateRefreshToken(payload dto.TokenPayload) (string, error) {
	claims := jwt.MapClaims{
		"userId":   payload.UserID.String(),
		"role":     payload.Role,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
		"iat":      time.Now().Unix(),
	}

	if payload.TenantID != nil {
		claims["tenantId"] = payload.TenantID.String()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GlobalConfig.JWTSecret))
}

func VerifyToken(tokenString string) (*dto.TokenPayload, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GlobalConfig.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	_ = claims // Temporarily suppressing unused warning for partial implementation

	// Porting payload back to struct
	// Note: Proper error handling for parsing UUIDs should be here
	return &dto.TokenPayload{
		// ... mapping claims ...
	}, nil
}
