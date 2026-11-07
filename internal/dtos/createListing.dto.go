package dtos

import (
	"fmt"
	"strings"
	"time"
)

// Request dto
type CreateListingRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	City        string `json:"city"`
}

type ValidationError struct {
	Field string
	Msg   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s : %s", e.Field, e.Msg)
}

func (req CreateListingRequest) Validate() error {
	return validateListingFields(req.Title, req.Description, req.City, req.Price)
}

func validateListingFields(title, description, city string, price int64) error {
	if strings.TrimSpace(title) == "" {
		return &ValidationError{Field: "title", Msg: "must not be empty"}
	}
	if len(title) > 200 {
		return &ValidationError{Field: "title", Msg: "must be at most 200 characters"}
	}
	if strings.TrimSpace(description) == "" {
		return &ValidationError{Field: "description", Msg: "must not be empty"}
	}
	if len(description) > 5000 {
		return &ValidationError{Field: "description", Msg: "must be at most 5000 characters"}
	}
	if price <= 0 {
		return &ValidationError{Field: "price", Msg: "must be greater than zero"}
	}
	if strings.TrimSpace(city) == "" {
		return &ValidationError{Field: "city", Msg: "must not be empty"}
	}

	return nil
}

// but for dto validation in go generally they use validator library 

// https://github.com/go-playground/validator

// But why we wont use it because it add dependency on external library and also we cant get compile time error this validation happens at run time 

// Response dto
type CreateListingResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}
