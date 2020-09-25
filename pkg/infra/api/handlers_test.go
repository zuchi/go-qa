package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/zuchi/go-qa/pkg/domain"
	"github.com/zuchi/go-qa/pkg/domain/question"
)

func TestServer_getListQuestion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := NewServer(context.Background(), nil)
	r := s.configHandlers(QuestionServiceDummy{})

	rNok := s.configHandlers(QuestionServiceDummyWithError{})

	tests := []struct {
		name          string
		method        string
		target        string
		engine        *gin.Engine
		length        int
		statusCode    int
		expectedError bool
	}{
		{
			name:          "should retrieve question list",
			method:        "GET",
			target:        "/question/",
			engine:        r,
			length:        1,
			statusCode:    http.StatusOK,
			expectedError: false,
		},
		{
			name:          "should get internal server error",
			method:        "GET",
			target:        "/question/",
			engine:        rNok,
			length:        0,
			statusCode:    http.StatusInternalServerError,
			expectedError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			request := httptest.NewRequest(test.method, test.target, nil)

			test.engine.ServeHTTP(w, request)

			assert.Equalf(t, test.statusCode, w.Code, "%s expected statusCode %d, but receive %d", test.name, test.statusCode, w.Code)

			if !test.expectedError {
				qb, err := ioutil.ReadAll(w.Body)
				assert.Nilf(t, err, "cannot parsed request result")

				var questions []*domain.Question
				err = json.Unmarshal(qb, &questions)
				assert.Nilf(t, err, "cannot parsed to slice of questions")

				assert.Lenf(t, questions, test.length, "%s expected 1, but receive %d", len(questions))
			}
		})
	}
}

func TestServer_postQuestion(t *testing.T) {
	gin.SetMode(gin.TestMode)

	dummyOk := QuestionServiceDummy{}
	dummyNOk := QuestionServiceDummyWithError{}

	q := `
	{
       "question": "What is the best language in the world?",
       "user": "Zuchi, Jederson",
       "createdDate": "0000-12-31T20:53:32-03:06"
    }`

	qi := `
	{
       "user": "Zuchi, Jederson",
       "createdDate": "0000-12-31T20:53:32-03:06"
    }`

	tests := []struct {
		name          string
		method        string
		target        string
		body          string
		questionPort  question.ServicePort
		statusCode    int
		message       string
		expectedError bool
	}{
		{
			name:          "should saved question",
			method:        "POST",
			target:        "/question/",
			body:          q,
			questionPort:  dummyOk,
			statusCode:    http.StatusCreated,
			message:       fmt.Sprintf("%s %s", "saved object id:", "SavedID"),
			expectedError: false,
		},
		{
			name:         "shouldn't not saved because question is invalid payload",
			method:       "POST",
			target:       "/question/",
			body:         qi,
			questionPort: dummyNOk,
			statusCode:   http.StatusBadRequest,
			message:      fmt.Sprintf("%s", "invalid payload"),
		},
		{
			name:         "shouldn't not saved because servicePort Retrieve error",
			method:       "POST",
			target:       "/question/",
			body:         q,
			questionPort: dummyNOk,
			statusCode:   http.StatusInternalServerError,
			message:      fmt.Sprintf("%s", "cannot save domain into repository"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			reader := strings.NewReader(test.body)
			request := httptest.NewRequest(test.method, test.target, reader)
			request.Header.Add("Content-Type", "application/json")

			s := NewServer(context.Background(), nil)
			r := s.configHandlers(test.questionPort)
			r.ServeHTTP(w, request)

			assert.Equalf(t, test.statusCode, w.Code, "%s expected statusCode %d, but receive %d", test.name, test.statusCode, w.Code)

			body, err := ioutil.ReadAll(w.Body)
			assert.Nilf(t, err, "%s cannot parsed request result", test.name)
			un := make(map[string]interface{})

			err = json.Unmarshal(body, &un)
			assert.Nilf(t, err, "%s cannot parsed to JSON", test.name)
			v, ok := un["message"]
			assert.True(t, ok, "%s expected message in json return didn't receive")
			assert.Equalf(t, test.message, v.(string), "%s expected %s, but receipt  %s", test.name, test.message, v.(string))
		})
	}
}

func TestServer_putQuestion(t *testing.T) {
	gin.SetMode(gin.TestMode)

	dummyOk := QuestionServiceDummy{}
	dummyNOk := QuestionServiceDummyWithError{}

	q := `
	{
       "question": "What is the best language in the world?",
       "user": "Zuchi, Jederson",
       "createdDate": "0000-12-31T20:53:32-03:06",
	   "answer": "Golang"
    }`

	qi := `
	{
	  "user": "Zuchi, Jederson",
	  "createdDate": "0000-12-31T20:53:32-03:06"
	}`

	ipMessage := fmt.Sprintf("%s", "invalid payload")
	qcbuMessage := fmt.Sprintf("%s", "question cannot be updated")

	tests := []struct {
		name          string
		method        string
		target        string
		body          string
		questionPort  question.ServicePort
		statusCode    int
		message       *string
		expectedError bool
	}{
		{
			name:          "should update question",
			method:        "PUT",
			target:        "/question/updateId",
			body:          q,
			questionPort:  dummyOk,
			statusCode:    http.StatusNoContent,
			message:       nil,
			expectedError: false,
		},
		{
			name:         "shouldn't not update because question is invalid payload",
			method:       "PUT",
			target:       "/question/updateId",
			body:         qi,
			questionPort: dummyNOk,
			statusCode:   http.StatusBadRequest,
			message:      &ipMessage,
		},
		{
			name:         "shouldn't not updated because servicePort Retrieve error",
			method:       "PUT",
			target:       "/question/updateId",
			body:         q,
			questionPort: dummyNOk,
			statusCode:   http.StatusInternalServerError,
			message:      &qcbuMessage,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			reader := strings.NewReader(test.body)
			request := httptest.NewRequest(test.method, test.target, reader)
			request.Header.Add("Content-Type", "application/json")

			s := NewServer(context.Background(), nil)
			r := s.configHandlers(test.questionPort)
			r.ServeHTTP(w, request)

			assert.Equalf(t, test.statusCode, w.Code, "%s expected statusCode %d, but receive %d", test.name, test.statusCode, w.Code)

			body, err := ioutil.ReadAll(w.Body)
			if test.statusCode != http.StatusNoContent {
				assert.Nilf(t, err, "%s cannot parsed request result", test.name)
				un := make(map[string]interface{})
				err = json.Unmarshal(body, &un)
				assert.Nilf(t, err, "%s cannot parsed to JSON", test.name)
				v, ok := un["message"]
				assert.True(t, ok, "%s expected message in json return didn't receive")
				assert.EqualValuesf(t, *test.message, v.(string), "%s expected %s, but receipt  %s", test.name, test.message, v.(string))
			}
		})
	}
}

// Dummy Mocks
type QuestionServiceDummy struct {
}

func (q QuestionServiceDummy) GetAllQuestion() ([]*domain.Question, error) {
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

func (q QuestionServiceDummy) SaveQuestion(_ *domain.Question) (id *string, err error) {
	ret := "SavedID"
	return &ret, nil
}

func (q QuestionServiceDummy) UpdateQuestion(_ string, _ *domain.Question) error {
	return nil
}

type QuestionServiceDummyWithError struct {
}

func (q QuestionServiceDummyWithError) GetAllQuestion() ([]*domain.Question, error) {
	return nil, errors.New("some error")
}

func (q QuestionServiceDummyWithError) SaveQuestion(question *domain.Question) (id *string, err error) {
	return nil, errors.New("some saved error")
}

func (q QuestionServiceDummyWithError) UpdateQuestion(_ string, _ *domain.Question) error {
	return errors.New("error to update")
}
