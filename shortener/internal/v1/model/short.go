package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	Short struct {
		ID       primitive.ObjectID `bson:"_id"`
		UserID   string             `bson:"user_id"`
		FullURL  string             `bson:"full_url"`
		ShortURL string             `bson:"short_url"`
		Visited  int64              `bson:"visited"`
	}

	CreateShortRequest struct {
		UserID   string `json:"user_id"`
		FullURL  string `json:"full_url"`
		ShortURL string `json:"short_url"`
	}
)
