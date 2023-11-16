package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ini buat sturck usernya
type ProductAll struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        *string            `json:"name" validate:"required,min=2,max=100"`        //validasi required yang di perlukan, min 2 karakter, max 100
	Description *string            `json:"description" validate:"required,min=2,max=100"` //validasi required yang di perlukan, min 2 karakter, max 100
	Price       *string            `json:"price" validate:"required"`                     //validasi required yang di perlukan, min 2 karakter
	Stock       *string            `json:"stock" validate:"required"`                     //validasi required yang di perlukan email wajib
	Size        *string            `json:"size" validate:"required"`
	Image       *string            `json:"image" validate:"required"`
}
