package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FaqAll struct {
	ID         primitive.ObjectID `bson:"_id"`
	Question   *string            `json:"question" validate:"required,min=2,max=100"`
	Answer     *string            `json:"answer"`
	Created_at time.Time          `json:"created_at"`
	Is_publish bool               `json:"is_publish"`
}
