package user

import "time"

type User struct {
	ID             string
	FirstName      string
	LastName       string
	Email          string
	HashedPassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
