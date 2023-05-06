package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	// User consist data of users
	User struct {
		ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
		FullName string             `bson:"fullname,omitempty" json:"full_name"`
		Email    string             `bson:"email,omitempty" json:"email"`
	}
	UserShorts struct {
		ID       string `json:"id"`
		FullURL  string `json:"full_url"`
		ShortURL string `json:"short_url"`
		Visited  int64  `json:"visited"`
	}
)
