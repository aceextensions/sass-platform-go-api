package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// QuestionType defines the format of the question
type QuestionType string

const (
	QuestionTypeMCQ       QuestionType = "mcq"
	QuestionTypeTrueFalse QuestionType = "true_false"
	QuestionTypeMatching  QuestionType = "matching"
)

// Question represents a test or quiz question
type Question struct {
	ID          uuid.UUID       `json:"id"`
	TenantID    uuid.UUID       `json:"tenantId"`
	SubjectID   uuid.UUID       `json:"subjectId" validate:"required"`
	Type        QuestionType    `json:"type" validate:"required,oneof=mcq true_false matching"`
	Category    string          `json:"category"`
	Difficulty  string          `json:"difficulty"` // e.g., beginner, intermediate, advanced
	Content     string          `json:"question" validate:"required"`
	Options     json.RawMessage `json:"options"` // Can be string array, object mappings, etc.
	Answer      string          `json:"answer" validate:"required"`
	Reasoning   string          `json:"reasoning"`
	Tags        []string        `json:"tags"`
	IsActive    bool            `json:"isActive"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// NewQuestion creates a new Question
func NewQuestion(tenantID, subjectID uuid.UUID, qType QuestionType, content string, answer string) *Question {
	now := time.Now()
	return &Question{
		ID:         uuid.New(),
		TenantID:   tenantID,
		SubjectID:  subjectID,
		Type:       qType,
		Content:    content,
		Answer:     answer,
		Tags:       make([]string, 0),
		IsActive:   true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// Common helper functions for options depending on type
func (q *Question) SetMCQOptions(options []string) error {
	b, err := json.Marshal(options)
	if err != nil {
		return err
	}
	q.Options = b
	return nil
}

func (q *Question) SetMatchingOptions(pairs map[string]string) error {
	b, err := json.Marshal(pairs)
	if err != nil {
		return err
	}
	q.Options = b
	return nil
}

func (q *Question) GetMCQOptions() ([]string, error) {
	var opts []string
	if len(q.Options) == 0 {
		return opts, nil
	}
	err := json.Unmarshal(q.Options, &opts)
	return opts, err
}
