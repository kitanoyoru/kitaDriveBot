package bcrypt

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/kitanoyoru/kitaDriveBot/libs/hasher"
)

func NewPasswordHasher(cost int) hasher.Hasher {
	return &bcryptHasher{cost}
}

type bcryptHasher struct {
	cost int
}

func (h *bcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	return string(bytes), err
}

func (h *bcryptHasher) CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
