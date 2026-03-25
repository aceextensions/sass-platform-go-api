package handler

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all quiz endpoints under a given router group
func RegisterRoutes(g *echo.Group, handler *QuizHandler) {
	quiz := g.Group("/quiz")

	// Subjects APIs
	quiz.GET("/subjects", handler.ListSubjects)
	quiz.POST("/subjects", handler.CreateSubject)
	quiz.GET("/subjects/:id", handler.GetSubject)

	// Questions APIs
	quiz.POST("/questions", handler.CreateQuestion)
	quiz.POST("/mock", handler.GenerateMockQuiz)
	quiz.GET("/questions", handler.GetQuestionsList)
}
