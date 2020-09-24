package question

import (
	"github.com/zuchi/go-qa/pkg/domain"
)

type Repository interface {
	GetAllQuestions() ([]*domain.Question, error)
	SaveQuestion(q *domain.Question) (*string, error)
	UpdateQuestions(id string, q *domain.Question) error
}
