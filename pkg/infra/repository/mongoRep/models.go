package mongoRep

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/zuchi/go-qa/pkg/domain"
)

type Question struct {
	Id          primitive.ObjectID  `bson:"_id"`
	Question    string              `bson:"question"`
	User        string              `bson:"user"`
	CreatedDate primitive.DateTime  `bson:"createdDate"`
	UpdateDate  *primitive.DateTime `bson:"updatedDate"`
	Answer      *string             `bson:"answer"`
}

func (q *Question) fromDomain(qd *domain.Question) error {
	if qd == nil {
		return nil
	}

	q.Question = qd.Question
	q.Answer = qd.Answer
	q.User = qd.User
	q.CreatedDate = primitive.NewDateTimeFromTime(qd.CreatedDate)

	if qd.UpdateDate != nil {
		ud := primitive.NewDateTimeFromTime(*qd.UpdateDate)
		q.UpdateDate = &ud
	}

	if qd.Id != nil {
		oId, err := primitive.ObjectIDFromHex(*qd.Id)
		if err != nil {
			return err
		}

		q.Id = oId
	}
	return nil
}

func (q *Question) toDomain() *domain.Question {
	dq := new(domain.Question)

	u := q.Id.Hex()
	dq.Id = &u

	dq.Answer = q.Answer
	dq.CreatedDate = q.CreatedDate.Time()
	dq.Question = q.Question
	dq.User = q.User
	if q.UpdateDate != nil {
		ut := q.UpdateDate.Time()
		dq.UpdateDate = &ut
	}

	return dq
}
