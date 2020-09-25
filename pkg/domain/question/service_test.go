package question

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zuchi/go-qa/pkg/domain"
)

func TestService_GetAllQuestion(t *testing.T) {
	dr := &DummyQuestionRepository{}
	dre := &DummyQuestionRepositoryError{}

	tests := []struct {
		name      string
		qr        Repository
		wantError bool
	}{
		{
			name:      "should retrieve questions",
			qr:        dr,
			wantError: false,
		},
		{
			name:      "should not retrieve questions",
			qr:        dre,
			wantError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := NewService(test.qr)
			question, err := s.GetAllQuestion()

			if !test.wantError {
				assert.Nilf(t, err, "%s expected error nil, but receipt %v", test.name, err)
				assert.Lenf(t, question, 1, "%s expected length 1, but receipt %d", test.name, len(question))
				assert.EqualValuesf(t, "Question 1", question[0].Question, "%s, expected %s, but receipt %s", test.name, "Question 1", question[0].Question)
			} else {
				assert.NotNilf(t, err, "%s expected error not nil, but receipt nil", test.name)
				assert.Nilf(t, question, "%s expected question is nil, but receipt %v", test.name, question)
			}
		})
	}
}

func TestService_SaveQuestion(t *testing.T) {
	dr := &DummyQuestionRepository{}
	dre := &DummyQuestionRepositoryError{}

	q := &domain.Question{
		Question: "some questions",
		User:     "someone",
	}

	tests := []struct {
		name      string
		q         *domain.Question
		qr        Repository
		wantError bool
	}{
		{
			name:      "should save questions",
			q:         q,
			qr:        dr,
			wantError: false,
		},
		{
			name:      "should not save questions",
			q:         q,
			qr:        dre,
			wantError: true,
		},
		{
			name:      "should not save questions because question is nil",
			q:         nil,
			qr:        dre,
			wantError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := NewService(test.qr)
			id, err := s.SaveQuestion(test.q)

			if !test.wantError {
				assert.Nilf(t, err, "%s expected error nil, but receipt %v", test.name, err)
				assert.EqualValuesf(t, "StoredID", *id, "%s expected StoredID, but receipt %s", test.name, *id)
			} else {
				assert.NotNilf(t, err, "%s expected error not nil, but receipt nil", test.name)
				assert.Nilf(t, id, "%s expected id equal nil, but receipt %v", test.name, id)
			}
		})
	}
}

func TestService_UpdateQuestion(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

type DummyQuestionRepository struct {
}

func (d *DummyQuestionRepository) GetAllQuestions() ([]*domain.Question, error) {
	var values []*domain.Question
	id := "id"
	ud := time.Now().Add(1 * time.Second)
	values = append(values, &domain.Question{
		Id:          &id,
		Question:    "Question 1",
		User:        "User 1",
		CreatedDate: time.Now(),
		UpdateDate:  &ud,
		Answer:      nil,
	})
	return values, nil
}

func (d *DummyQuestionRepository) SaveQuestion(q *domain.Question) (*string, error) {
	id := "StoredID"
	return &id, nil
}

func (d *DummyQuestionRepository) UpdateQuestions(_ string, _ *domain.Question) error {
	return nil
}

type DummyQuestionRepositoryError struct {
}

func (d DummyQuestionRepositoryError) GetAllQuestions() ([]*domain.Question, error) {
	return nil, errors.New("some error here")
}

func (d DummyQuestionRepositoryError) SaveQuestion(_ *domain.Question) (*string, error) {
	return nil, errors.New("some saved error")
}

func (d DummyQuestionRepositoryError) UpdateQuestions(_ string, _ *domain.Question) error {
	return errors.New("some error here")
}
