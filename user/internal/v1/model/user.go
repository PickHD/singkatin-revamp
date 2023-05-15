package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	// User consist data of users
	User struct {
		ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
		FullName string             `bson:"fullname,omitempty" json:"full_name"`
		Email    string             `bson:"email,omitempty" json:"email"`
	}

	// UserShorts consist data of user shorts
	UserShorts struct {
		ID       string `json:"id"`
		FullURL  string `json:"full_url"`
		ShortURL string `json:"short_url"`
		Visited  int64  `json:"visited"`
	}

	// GenerateShortUserRequest consist request data generate short users
	GenerateShortUserRequest struct {
		FullURL string `json:"full_url"`
	}

	// GenerateShortUserResponse consist response data when success generate short users
	GenerateShortUserResponse struct {
		ShortURL string `json:"short_url"`
	}

	// GenerateShortUserMessage consist message short users to publish
	GenerateShortUserMessage struct {
		FullURL  string `json:"full_url"`
		ShortURL string `json:"short_url"`
		UserID   string `json:"user_id"`
	}

	EditProfileRequest struct {
		FullName string `json:"full_name"`
	}
)
