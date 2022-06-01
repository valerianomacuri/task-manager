package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	User struct {
		Id           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
		FirstName    string             `json:"firstname"`
		LastName     string             `json:"lastname"`
		Email        string             `json:"email"`
		Password     string             `json:"password,omitempty"`
		HashPassword []byte             `json:"hashpassword,omitempty"`
	}
	Task struct {
		Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
		CreatedBy   string             `json:"createdby"`
		Name        string             `json:"name"`
		Description string             `json:"description"`
		CreatedOn   time.Time          `json:"createdon,omitempty"`
		Due         time.Time          `json:"due,omitempty"`
		Status      string             `json:"status,omitempty"`
		Tags        []string           `json:"tags,omitempty"`
	}
	TaskNote struct {
		Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
		TaskId      primitive.ObjectID `json:"taskid"`
		Description string             `json:"description"`
		CreatedOn   time.Time          `json:"createdon,omitempty"`
	}
)
