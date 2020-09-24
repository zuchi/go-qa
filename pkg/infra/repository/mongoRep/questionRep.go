package mongoRep

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/zuchi/go-qa/pkg/domain"
	"github.com/zuchi/go-qa/pkg/domain/question"
)

type QuestionRepository struct {
	ctx    context.Context
	qc     *mongo.Collection
	logCtx *log.Entry
}

func (m *QuestionRepository) GetAllQuestions() ([]*domain.Question, error) {
	logF := m.logCtx.WithField("function", "GetAllQuestions")
	logF.Info("Getting all Questions")

	cursor, err := m.qc.Find(m.ctx, bson.M{})

	defer func() {
		err := cursor.Close(m.ctx)
		if err != nil {
			logF.Warnf("cannot close cursor: %v", err)
		}
	}()

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var questions []Question
	err = cursor.All(m.ctx, &questions)
	if err != nil {
		return nil, err
	}

	qDomains := make([]*domain.Question, 0)
	for i := range questions {
		qDomains = append(qDomains, questions[i].toDomain())
	}

	return qDomains, nil
}

func (m *QuestionRepository) SaveQuestion(qd *domain.Question) (*string, error) {
	logF := m.logCtx.WithField("function", "SaveQuestion")
	logF.Info("Saving new question")
	if qd == nil {
		return nil, errors.New("there is no question to be saved")
	}

	var q Question
	err := q.fromDomain(qd)
	if err != nil {
		return nil, err
	}

	q.Id = primitive.NewObjectID()

	one, err := m.qc.InsertOne(m.ctx, q)
	if err != nil {
		return nil, err
	}

	oId, ok := one.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("cannot get inserted Id")
	}

	id := oId.Hex()
	return &id, nil
}

func (m *QuestionRepository) UpdateQuestions(id string, qd *domain.Question) error {
	logF := m.logCtx.WithField("function", "SaveQuestion")
	logF.Info("Updating new question")
	if qd == nil {
		return errors.New("there is no question to be updated")
	}

	var q Question
	err := q.fromDomain(qd)
	if err != nil {
		return err
	}

	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	q.Id = oId

	_, err = m.qc.UpdateOne(m.ctx, bson.M{"_id": oId}, bson.M{"$set": q})
	return err
}

func NewQuestionRepository(ctx context.Context, qc *mongo.Collection) question.Repository {
	logCtx := log.WithFields(log.Fields{"component": "mongoRep.Repository"})
	return &QuestionRepository{
		ctx:    ctx,
		qc:     qc,
		logCtx: logCtx,
	}
}
