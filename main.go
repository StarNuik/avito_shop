package main

import (
	"github.com/avito_shop/internal/setup"
)

// testing:
// go test -v -coverpkg=./... -coverprofile=./all.cov ./...
// go tool cover -html=all.cov

func main() {
	router := setup.Router()
	_ = router.Run()
}
