package model

import "regexp"

type (
	// RegisterRequest consist request data for registering as users
	RegisterRequest struct {
		FullName string `json:"fullname"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// RegisterResponse consist response of success registering as users
	RegisterResponse struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}
)

var (
	IsValidEmail, _ = regexp.Compile(`^(?P<name>[a-zA-Z0-9.!#$%&'*+/=?^_ \x60{|}~-]+)@(?P<domain>[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*)$`)
)
