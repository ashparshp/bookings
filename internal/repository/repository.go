package repository

import "github.com/ashparshp/bookings/internal/models"
type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end string, roomID int) (bool, error)
}

