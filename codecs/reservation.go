package codecs

import (
	"easybook/models"
	"easybook/types"
)

// BookingReserveRoomsRequest is a struct for reserving rooms.
type BookingReserveRoomsRequest struct {
	GuestID         int        `json:"guestId" validate:"required"`
	StartDate       types.Date `json:"startDate" validate:"required"`
	EndDate         types.Date `json:"endDate" validate:"required"`
	DiscountPercent float32    `json:"discountPercent,omitempty"`
	TotalPrice      float32    `json:"totalPrice" validate:"required"`
	Rooms           []int      `json:"rooms" validate:"required"`
}

// BookingReserveRoomsResponse is a struct for reserving rooms.
type BookingReserveRoomsResponse struct {
	CommonResponse
	Reservation *models.Reservation `json:"reservation,omitempty"`
}
