package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Short struct {
	ID       primitive.ObjectID `bson:"_id"`
	UserID   primitive.ObjectID `bson:"user_id"`
	FullURL  string             `bson:"full_url"`
	ShortURL string             `bson:"short_url"`
	Visited  int64              `bson:"visited"`
}
