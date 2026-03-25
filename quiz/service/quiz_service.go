package service

import (
	"context"

	"github.com/aceextension/quiz/domain"
	"github.com/aceextension/quiz/repository"
	"github.com/google/uuid"
)

type QuizService interface {
	// Subject methods
	CreateSubject(ctx context.Context, tenantID uuid.UUID, name, description string) (*domain.Subject, error)
	GetSubject(ctx context.Context, id uuid.UUID) (*domain.Subject, error)
	ListSubjects(ctx context.Context) ([]*domain.Subject, error)

	// Question methods
	CreateQuestion(ctx context.Context, tenantID uuid.UUID, q *domain.Question) error
	GetQuestions(ctx context.Context, filter repository.QuestionFilter) ([]*domain.Question, int, error)
	GenerateMockQuiz(ctx context.Context, filter repository.QuestionFilter) ([]*domain.Question, error)
}

type quizService struct {
	subjectRepo  repository.SubjectRepository
	questionRepo repository.QuestionRepository
}

func NewQuizService(subjectRepo repository.SubjectRepository, questionRepo repository.QuestionRepository) QuizService {
	return &quizService{
		subjectRepo:  subjectRepo,
		questionRepo: questionRepo,
	}
}

// Subject methods

func (s *quizService) CreateSubject(ctx context.Context, tenantID uuid.UUID, name, description string) (*domain.Subject, error) {
	subject := domain.NewSubject(tenantID, name, description)
	if err := s.subjectRepo.Create(ctx, subject); err != nil {
		return nil, err
	}
	return subject, nil
}

func (s *quizService) GetSubject(ctx context.Context, id uuid.UUID) (*domain.Subject, error) {
	return s.subjectRepo.GetByID(ctx, id)
}

func (s *quizService) ListSubjects(ctx context.Context) ([]*domain.Subject, error) {
	return s.subjectRepo.List(ctx)
}

// Question methods

func (s *quizService) CreateQuestion(ctx context.Context, tenantID uuid.UUID, q *domain.Question) error {
	q.TenantID = tenantID
	return s.questionRepo.Create(ctx, q)
}

func (s *quizService) GetQuestions(ctx context.Context, filter repository.QuestionFilter) ([]*domain.Question, int, error) {
	return s.questionRepo.List(ctx, filter)
}

func (s *quizService) GenerateMockQuiz(ctx context.Context, filter repository.QuestionFilter) ([]*domain.Question, error) {
	return s.questionRepo.GetRandomQuestions(ctx, filter)
}
