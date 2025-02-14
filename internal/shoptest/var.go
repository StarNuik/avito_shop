package shoptest

const (
	HostUrl        = "http://localhost:8080"
	DefaultBalance = int64(1000)
)

type user struct {
	Username string
	Password string
}

var Users = []user{
	{Username: "user#0", Password: "user#0"},
	{Username: "user#1", Password: "user#1"},
	{Username: "user#2", Password: "user#2"},
}
