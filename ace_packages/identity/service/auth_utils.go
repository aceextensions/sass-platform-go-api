package service

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aceextension/core/config"
	"github.com/aceextension/identity/dto"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

const (
	Argon2Memory      = 65536
	Argon2Iterations  = 3
	Argon2Parallelism = 4
	Argon2SaltLen     = 16
	Argon2KeyLen      = 32
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, Argon2SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, Argon2Iterations, Argon2Memory, Argon2Parallelism, Argon2KeyLen)

	// Format: $argon2id$v=19$m=65536,t=3,p=4$salt$hash
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, Argon2Memory, Argon2Iterations, Argon2Parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func ComparePassword(password, encodedHash string) bool {
	// Format: $argon2id$v=19$m=65536,t=3,p=4$salt$hash
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false
	}

	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil || version != argon2.Version {
		return false
	}

	var memory uint32
	var iterations uint32
	var parallelism uint8
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}

	keyLen := uint32(len(decodedHash))

	comparisonHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLen)

	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1
}

func GenerateAccessToken(payload dto.TokenPayload) (string, error) {
	claims := jwt.MapClaims{
		"userId": payload.UserID.String(),
		"role":   payload.Role,
		"exp":    time.Now().Add(time.Hour * 1).Unix(), // 1 hour
		"iat":    time.Now().Unix(),
	}

	if payload.TenantID != nil {
		claims["tenantId"] = payload.TenantID.String()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GlobalConfig.JWTSecret))
}

func GenerateRefreshToken(payload dto.TokenPayload) (string, error) {
	claims := jwt.MapClaims{
		"userId": payload.UserID.String(),
		"role":   payload.Role,
		"exp":    time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
		"iat":    time.Now().Unix(),
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
