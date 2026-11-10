package dtos

import (
	"net/mail"
)

type SigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req SigninRequest) Validate() error {
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return &ValidationError{Field: "email", Msg: "must be a valid email address"}
	}

	// if len(req.Password) < 8 {
	// 	return &ValidationError{Field: "password", Msg: "must be at least 8 characters"}
	// }

	return nil
}

type SigninResponse struct{
	Token string
	ExpiresIn int
}