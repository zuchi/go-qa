package domain

import (
	"time"
)

type Question struct {
	Id          *string    `json:"id"`
	Question    string     `json:"question" binding:"required"`
	User        string     `json:"user"`
	CreatedDate time.Time  `json:"createdDate"`
	UpdateDate  *time.Time `json:"updatedDate"`
	Answer      *string    `json:"answer"`
}
