package hasher

type Hasher interface {
	Hash(password string) (string, error)
	CheckPasswordHash(password, hashedPassword string) bool
}

type Generator interface {
	Generate(n int) (string, error)
}
