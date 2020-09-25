package api

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/zuchi/go-qa/pkg/domain/question"
	"github.com/zuchi/go-qa/pkg/infra/repository/mongoRep"
)

type Server struct {
	ctx    context.Context
	srv    *http.Server
	logCtx *log.Entry
	client *mongo.Client
}

func NewServer(ctx context.Context, mongo *mongo.Client) *Server {
	l := log.WithFields(log.Fields{"component": "server"})

	return &Server{
		ctx:    ctx,
		logCtx: l,
		client: mongo,
	}
}

func (s *Server) configHandlers(qs question.ServicePort) *gin.Engine {
	s.logCtx.Infof("registering routes")
	r := gin.New()

	r.Use(ginrus.Ginrus(s.logCtx, time.RFC3339, true))
	r.Use(gin.Recovery())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Pong"})
	})

	g := r.Group("/question/")
	g.GET("", s.getListQuestion(qs))
	g.POST("", s.postQuestion(qs))
	g.PUT(":id", s.putQuestion(qs))

	return r
}

func (s *Server) getQuestionService() question.ServicePort {
	db := os.Getenv("MONGO_COLLECTION")
	questions := s.client.Database(db).Collection("questions")
	repository := mongoRep.NewQuestionRepository(s.ctx, questions)

	return question.NewService(repository)
}

func (s *Server) Initialize(Addr string) {

	s.srv = &http.Server{
		Addr:    Addr,
		Handler: s.configHandlers(s.getQuestionService()),
	}

	s.logCtx.Infof("starting server addr: %s", s.srv.Addr)

	go func() {
		err := s.srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			s.logCtx.Fatalf("%v", err)
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logCtx.Info("shutdown serving...")
	return s.srv.Shutdown(ctx)
}
