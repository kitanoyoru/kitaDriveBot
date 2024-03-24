package bcrypt

import "testing"

func TestLol(t *testing.T) {
	h := NewPasswordHasher(5)
	t.Log(h.Hash("111111"))
}
