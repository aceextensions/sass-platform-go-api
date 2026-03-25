package handler

import (
	"net/http"
	"strconv"

	"github.com/aceextension/core/apperrors"
	"github.com/aceextension/quiz/domain"
	"github.com/aceextension/quiz/repository"
	"github.com/aceextension/quiz/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type QuizHandler struct {
	service service.QuizService
}

func NewQuizHandler(svc service.QuizService) *QuizHandler {
	return &QuizHandler{service: svc}
}

// Subject Handlers

type CreateSubjectRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

func (h *QuizHandler) CreateSubject(c echo.Context) error {
	var req CreateSubjectRequest
	if err := c.Bind(&req); err != nil {
		return apperrors.NewAppError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := c.Validate(req); err != nil {
		return err // Let global error handler deal with validation errors natively
	}

	// For mock purposes we generate a dummy tenant ID
	tenantID := uuid.New()

	subject, err := h.service.CreateSubject(c.Request().Context(), tenantID, req.Name, req.Description)
	if err != nil {
		return apperrors.NewAppError(http.StatusInternalServerError, "Failed to create subject")
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    subject,
	})
}

func (h *QuizHandler) GetSubject(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return apperrors.NewAppError(http.StatusBadRequest, "Invalid subject ID")
	}

	subject, err := h.service.GetSubject(c.Request().Context(), id)
	if err != nil {
		return apperrors.NewAppError(http.StatusNotFound, "Subject not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    subject,
	})
}

func (h *QuizHandler) ListSubjects(c echo.Context) error {
	subjects, err := h.service.ListSubjects(c.Request().Context())
	if err != nil {
		return apperrors.NewAppError(http.StatusInternalServerError, "Failed to retrieve subjects")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    subjects,
	})
}

// Question Handlers

type GenerateMockQuizRequest struct {
	SubjectID  *string  `json:"subject_id"`
	Category   *string  `json:"category"`
	Difficulty *string  `json:"difficulty"`
	Limit      int      `json:"limit"`
	Types      []string `json:"types"`
}

func (h *QuizHandler) GenerateMockQuiz(c echo.Context) error {
	req := new(GenerateMockQuizRequest)
	if err := c.Bind(req); err != nil {
		return apperrors.NewAppError(http.StatusBadRequest, "Invalid request payload")
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10 // Default to 10
	}

	filter := repository.QuestionFilter{
		Limit: limit,
	}

	if req.SubjectID != nil && *req.SubjectID != "" {
		id, err := uuid.Parse(*req.SubjectID)
		if err == nil {
			filter.SubjectID = &id
		}
	}

	if req.Category != nil && *req.Category != "" {
		filter.Category = req.Category
	}

	if req.Difficulty != nil && *req.Difficulty != "" {
		filter.Difficulty = req.Difficulty
	}

	// For simplicity, we just look at the first type if passed
	if len(req.Types) > 0 && req.Types[0] != "" {
		t := domain.QuestionType(req.Types[0])
		filter.Type = &t
	}

	questions, err := h.service.GenerateMockQuiz(c.Request().Context(), filter)
	if err != nil {
		return apperrors.NewAppError(http.StatusInternalServerError, "Failed to generate mock quiz")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    questions,
		"message": "Successfully generated mock questions",
	})
}

type CreateQuestionRequest struct {
	SubjectID  string   `json:"subject_id" validate:"required"`
	Type       string   `json:"type" validate:"required"`
	Question   string   `json:"question" validate:"required"`
	Answer     string   `json:"answer" validate:"required"`
	Options    []string `json:"options"`
	Category   string   `json:"category"`
	Difficulty string   `json:"difficulty"`
	Reasoning  string   `json:"reasoning"`
	Tags       []string `json:"tags"`
}

func (h *QuizHandler) CreateQuestion(c echo.Context) error {
	var req CreateQuestionRequest
	if err := c.Bind(&req); err != nil {
		return apperrors.NewAppError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	subjectID, err := uuid.Parse(req.SubjectID)
	if err != nil {
		return apperrors.NewAppError(http.StatusBadRequest, "Invalid subject ID")
	}

	// For mock purposes we generate a dummy tenant ID
	tenantID := uuid.New()

	q := domain.NewQuestion(tenantID, subjectID, domain.QuestionType(req.Type), req.Question, req.Answer)
	q.Category = req.Category
	q.Difficulty = req.Difficulty
	q.Reasoning = req.Reasoning
	q.Tags = req.Tags
	q.SetMCQOptions(req.Options)

	if err := h.service.CreateQuestion(c.Request().Context(), tenantID, q); err != nil {
		return apperrors.NewAppError(http.StatusInternalServerError, "Failed to create question")
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    q,
	})
}

func (h *QuizHandler) GetQuestionsList(c echo.Context) error {
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	page, _ := strconv.Atoi(pageStr)
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 20
	}

	offset := (page - 1) * limit

	filter := repository.QuestionFilter{
		Limit:  limit,
		Offset: offset,
	}

	// Read optional query params
	if sID := c.QueryParam("subject_id"); sID != "" {
		id, err := uuid.Parse(sID)
		if err == nil {
			filter.SubjectID = &id
		}
	}
	if cat := c.QueryParam("category"); cat != "" {
		filter.Category = &cat
	}
	if diff := c.QueryParam("difficulty"); diff != "" {
		filter.Difficulty = &diff
	}
	if t := c.QueryParam("type"); t != "" {
		qType := domain.QuestionType(t)
		filter.Type = &qType
	}

	questions, total, err := h.service.GetQuestions(c.Request().Context(), filter)
	if err != nil {
		return apperrors.NewAppError(http.StatusInternalServerError, "Failed to retrieve questions")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    questions,
		"meta": map[string]interface{}{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}
