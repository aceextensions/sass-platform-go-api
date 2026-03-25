package repository

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/aceextension/quiz/domain"
	"github.com/google/uuid"
)

type subjectMemoryRepository struct {
	mu       sync.RWMutex
	subjects map[uuid.UUID]*domain.Subject
}

type questionMemoryRepository struct {
	mu        sync.RWMutex
	questions map[uuid.UUID]*domain.Question
}

func NewMemoryRepository() (SubjectRepository, QuestionRepository) {
	sRepo := &subjectMemoryRepository{
		subjects: make(map[uuid.UUID]*domain.Subject),
	}
	qRepo := &questionMemoryRepository{
		questions: make(map[uuid.UUID]*domain.Question),
	}

	return sRepo, qRepo
}

// Ensure the types are correct structure for seeding
type seedQuestion struct {
	ID         int      `json:"id"`
	Domain     string   `json:"domain"`
	Type       string   `json:"type"` // Might not exist in exam.json, default to mcq
	Difficulty string   `json:"difficulty"`
	Question   string   `json:"question"`
	Options    []string `json:"options"`
	Answer     string   `json:"answer"`
	Reasoning  string   `json:"reasoning"`
	Tags       []string `json:"tags"`
}

func SeedFromFile(sRepo SubjectRepository, qRepo QuestionRepository, filePath string, defaultSubjectName string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var seeds []seedQuestion
	if err := json.Unmarshal(data, &seeds); err != nil {
		return err
	}

	tenantID := uuid.New() // Placeholder tenant
	subject := domain.NewSubject(tenantID, defaultSubjectName, "Auto-generated subject from seed file")

	if err := sRepo.Create(context.Background(), subject); err != nil {
		return err
	}

	for _, s := range seeds {
		qType := domain.QuestionTypeMCQ
		if s.Type == "true_false" {
			qType = domain.QuestionTypeTrueFalse
		} else if s.Type == "matching" {
			qType = domain.QuestionTypeMatching
		}

		q := domain.NewQuestion(tenantID, subject.ID, qType, s.Question, s.Answer)
		q.Category = s.Domain
		q.Difficulty = s.Difficulty
		q.Reasoning = s.Reasoning
		q.Tags = s.Tags
		
		q.SetMCQOptions(s.Options)
		if err := qRepo.Create(context.Background(), q); err != nil {
			return err
		}
	}

	return nil
}

// SubjectRepository implementation

func (r *subjectMemoryRepository) Create(ctx context.Context, subject *domain.Subject) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.subjects[subject.ID]; exists {
		return errors.New("subject already exists")
	}
	r.subjects[subject.ID] = subject
	return nil
}

func (r *subjectMemoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subject, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	subject, exists := r.subjects[id]
	if !exists {
		return nil, errors.New("subject not found")
	}
	return subject, nil
}

func (r *subjectMemoryRepository) GetByName(ctx context.Context, name string) (*domain.Subject, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, subject := range r.subjects {
		if subject.Name == name {
			return subject, nil
		}
	}
	return nil, errors.New("subject not found")
}

func (r *subjectMemoryRepository) List(ctx context.Context) ([]*domain.Subject, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*domain.Subject
	for _, subject := range r.subjects {
		list = append(list, subject)
	}
	return list, nil
}

func (r *subjectMemoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.subjects, id)
	return nil
}

// QuestionRepository implementation

func (r *questionMemoryRepository) Create(ctx context.Context, question *domain.Question) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.questions[question.ID]; exists {
		return errors.New("question already exists")
	}
	r.questions[question.ID] = question
	return nil
}

func (r *questionMemoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Question, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	question, exists := r.questions[id]
	if !exists {
		return nil, errors.New("question not found")
	}
	return question, nil
}

func (r *questionMemoryRepository) matchFilter(q *domain.Question, filter QuestionFilter) bool {
	if filter.SubjectID != nil && *filter.SubjectID != uuid.Nil && q.SubjectID != *filter.SubjectID {
		return false
	}
	if filter.Category != nil && q.Category != *filter.Category {
		return false
	}
	if filter.Type != nil && q.Type != *filter.Type {
		return false
	}
	if filter.Difficulty != nil && q.Difficulty != *filter.Difficulty {
		return false
	}
	return true
}

func (r *questionMemoryRepository) List(ctx context.Context, filter QuestionFilter) ([]*domain.Question, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*domain.Question
	for _, q := range r.questions {
		if r.matchFilter(q, filter) {
			filtered = append(filtered, q)
		}
	}

	total := len(filtered)
	
	// Apply pagination
	start := filter.Offset
	if start > total {
		start = total
	}
	end := start + filter.Limit
	if end > total || filter.Limit == 0 {
		end = total // if limit is 0, return all from offset
	}

	return filtered[start:end], total, nil
}

func (r *questionMemoryRepository) GetRandomQuestions(ctx context.Context, filter QuestionFilter) ([]*domain.Question, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*domain.Question
	for _, q := range r.questions {
		if r.matchFilter(q, filter) {
			filtered = append(filtered, q)
		}
	}

	limit := filter.Limit
	if limit == 0 || limit > len(filtered) {
		limit = len(filtered)
	}

	// Shuffle
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(filtered), func(i, j int) {
		filtered[i], filtered[j] = filtered[j], filtered[i]
	})

	return filtered[:limit], nil
}

func (r *questionMemoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.questions, id)
	return nil
}
