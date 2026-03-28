package quiz

import (
	"github.com/aceextension/core/extension"
	"github.com/aceextension/quiz/handler"
	"github.com/labstack/echo/v4"
)

type QuizModule struct {
}

func NewQuizModule() *QuizModule {
	return &QuizModule{}
}

func (m *QuizModule) Name() string {
	return "quiz"
}

func (m *QuizModule) Init() error {
	Init("")
	return nil
}

func (m *QuizModule) RegisterRoutes(e *echo.Echo, g *echo.Group) error {
	handler.RegisterRoutes(g, handler.NewQuizHandler(Service))
	return nil
}

func (m *QuizModule) RegisterEvents() error {
	return nil
}

func (m *QuizModule) RegisterPlugins(registry *extension.PluginRegistry) error {
	return nil
}
