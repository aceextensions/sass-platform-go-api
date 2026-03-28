package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/aceextension/user-management/domain"
	"github.com/google/uuid"
)

type SessionService interface {
	CreateSession(ctx context.Context, userID uuid.UUID, ip, userAgent string) (*domain.Session, error)
	VerifySession(ctx context.Context, token string) (*domain.Session, error)
	RevokeSession(ctx context.Context, sessionID uuid.UUID) error
}

type sessionService struct {
	// In a real implementation, this would use a repository
}

func NewSessionService() SessionService {
	return &sessionService{}
}

func (s *sessionService) CreateSession(ctx context.Context, userID uuid.UUID, ip, userAgent string) (*domain.Session, error) {
	token := uuid.New().String()
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	session := &domain.Session{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: tokenHash,
		IPAddress: ip,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	// TODO: Save session to database via repository
	fmt.Printf("📝 Session created for user %s (IP: %s)\n", userID, ip)

	return session, nil
}

func (s *sessionService) VerifySession(ctx context.Context, token string) (*domain.Session, error) {
	// TODO: Lookup session by token hash
	return nil, nil
}

func (s *sessionService) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	// TODO: Delete session
	return nil
}

// GenerateFingerprint creates a unique identifier for a device/browser combination
func GenerateFingerprint(ip, userAgent string) string {
	data := fmt.Sprintf("%s|%s", ip, userAgent)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
