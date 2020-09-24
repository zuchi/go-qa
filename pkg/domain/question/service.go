package question

import (
	"errors"
	"time"

	"github.com/zuchi/go-qa/pkg/domain"
)

type Service struct {
	questionRepo Repository
}

func NewService(questionRepo Repository) *Service {
	return &Service{
		questionRepo: questionRepo,
	}
}

func (s *Service) GetAllQuestion() ([]*domain.Question, error) {
	return s.questionRepo.GetAllQuestions()
}

func (s *Service) SaveQuestion(question *domain.Question) (id *string, err error) {
	if question == nil {
		return nil, errors.New("there is no question to be saved")
	}

	question.CreatedDate = time.Now()
	return s.questionRepo.SaveQuestion(question)
}

func (s *Service) UpdateQuestion(id string, question *domain.Question) error {
	if question == nil {
		return errors.New("there is no question to be updated")
	}
	question.Id = &id
	updatedDate := time.Now()
	question.UpdateDate = &updatedDate
	return s.questionRepo.UpdateQuestions(id, question)
}
