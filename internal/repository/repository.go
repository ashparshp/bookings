package repository

import "github.com/ashparshp/bookings/internal/models"
type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) error
}

