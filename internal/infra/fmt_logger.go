package infra

import "log"

type FmtLogger struct{}

var _ Logger = (*FmtLogger)(nil)

func (_ *FmtLogger) Log(message string) {
	log.Println(message)
}

func (_ *FmtLogger) LogError(err error) {
	log.Println(err)
}
