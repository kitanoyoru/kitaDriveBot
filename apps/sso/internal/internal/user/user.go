package user

import "time"

type User struct {
	ID             int64
	FirstName      string
	LastName       string
	Email          string
	HashedPassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
