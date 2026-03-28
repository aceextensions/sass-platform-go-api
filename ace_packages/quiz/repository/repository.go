package repository

import (
	"context"

	"github.com/aceextension/quiz/domain"
	"github.com/google/uuid"
)

// SubjectRepository defines operations for subjects
type SubjectRepository interface {
	Create(ctx context.Context, subject *domain.Subject) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Subject, error)
	GetByName(ctx context.Context, name string) (*domain.Subject, error)
	List(ctx context.Context) ([]*domain.Subject, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// QuestionFilter contains optional parameters to filter questions
type QuestionFilter struct {
	SubjectID  *uuid.UUID
	Category   *string
	Type       *domain.QuestionType
	Difficulty *string
	Limit      int
	Offset     int
}

// QuestionRepository defines operations for questions
type QuestionRepository interface {
	Create(ctx context.Context, question *domain.Question) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Question, error)
	List(ctx context.Context, filter QuestionFilter) ([]*domain.Question, int, error)
	GetRandomQuestions(ctx context.Context, filter QuestionFilter) ([]*domain.Question, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
