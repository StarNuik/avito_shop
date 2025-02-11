package infra

type Logger interface {
	Log(string)
	LogError(error)
}
