package codecs

import (
	"easybook/models"
)

// GuestPostRequest is a struct for creating or updating guest.
type GuestPostRequest struct {
	ID                int    `json:"id,omitempty"`
	FirstName         string `json:"firstName" validate:"required"`
	LastName          string `json:"lastName" validate:"required"`
	Email             string `json:"email" validate:"required"`
	Phone             string `json:"phone,omitempty"`
	Address           string `json:"address,omitempty"`
	Detail            string `json:"detail,omitempty"`
	Role              int8   `json:"role" validate:"required"`
	Password          string `json:"password,omitempty"`
	ConfirmedPassword string `json:"confirmedPassword,omitempty"`
}

// GuestPostResponse is a struct for return new created or updated guest.
type GuestPostResponse struct {
	CommonResponse
	Guest *models.Guest `json:"guest,omitempty"`
}
