package domain

type PasswordHasher interface {
	Hash(password string) (string, error)
	Same(inPassword string, storedHash string) bool
}
