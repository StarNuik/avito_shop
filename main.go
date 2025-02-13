package main

import (
	"github.com/avito_shop/internal/init"
)

func main() {
	router := init.Router()
	_ = router.Run()
}
