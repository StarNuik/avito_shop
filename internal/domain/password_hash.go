package domain

type PasswordHash interface {
	Hash(string) (string, error)
	Same(string, string) bool
}
