package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zuchi/go-qa/pkg/domain"
	"github.com/zuchi/go-qa/pkg/domain/question"
)

func (s *Server) getListQuestion(qs question.ServicePort) gin.HandlerFunc {
	return func(c *gin.Context) {
		questions, err := qs.GetAllQuestion()
		if err != nil {
			s.logCtx.WithError(err).Error("something went wrong to getAllQuestions")
			c.AbortWithStatusJSON(http.StatusInternalServerError,
				gin.H{"message": "something went wrong to get questions and answers"})
			return
		}

		c.JSON(http.StatusOK, questions)
	}
}

func (s *Server) postQuestion(qs question.ServicePort) gin.HandlerFunc {
	return func(c *gin.Context) {
		var q domain.Question
		if err := c.ShouldBind(&q); err != nil {
			s.logCtx.WithError(err).Error("something went wrong to parser domain.question")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid payload"})
			return
		}

		id, err := qs.SaveQuestion(&q)
		if err != nil {
			s.logCtx.Errorf("something went wrong to save domain.question: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "cannot save domain into repository"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "saved object id: " + *id})
	}
}

func (s *Server) putQuestion(qs question.ServicePort) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Params.ByName("id")

		var q domain.Question
		if err := c.ShouldBind(&q); err != nil {
			s.logCtx.WithError(err).Error("something went wrong to parser Question")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid payload"})
			return
		}

		err := qs.UpdateQuestion(id, &q)
		if err != nil {
			s.logCtx.WithError(err).Error("something went wrong to update question")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "question cannot be updated"})
			return
		}

		c.JSON(http.StatusNoContent, http.NoBody)
	}
}
