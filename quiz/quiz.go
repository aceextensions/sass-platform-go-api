package quiz

import (
	"log"

	"github.com/aceextension/quiz/repository"
	"github.com/aceextension/quiz/service"
)

var (
	Service      service.QuizService
	SubjectRepo  repository.SubjectRepository
	QuestionRepo repository.QuestionRepository
)

// Init initializes the quiz module and seeds data from the given filePath
// If no filePath is given, it operates empty
func Init(seedFile string) {
	subjectRepo, questionRepo := repository.NewMemoryRepository()

	SubjectRepo = subjectRepo
	QuestionRepo = questionRepo
	Service = service.NewQuizService(SubjectRepo, QuestionRepo)

	// Seed from file if provided
	if seedFile != "" {
		err := repository.SeedFromFile(subjectRepo, questionRepo, seedFile, "Docker Certification")
		if err != nil {
			log.Printf("Warning: failed to seed quiz repository from %s: %v\n", seedFile, err)
		} else {
			log.Println("Successfully seeded quiz questions from", seedFile)
		}
	}
}
