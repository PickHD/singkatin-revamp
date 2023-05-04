package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	// User consist data of users
	User struct {
		ID       primitive.ObjectID `bson:"_id,omitempty"`
		FullName string             `bson:"fullname,omitempty"`
		Email    string             `bson:"email,omitempty"`
		Password string             `bson:"password,omitempty"`
	}
)
