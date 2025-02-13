package shoptest

import (
	"github.com/avito_shop/internal/infra"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

func NewGinRoundtrip(router *gin.Engine) infra.RoundtripHandler {
	return func(req *http.Request) (*http.Response, error) {
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		return recorder.Result(), nil
	}
}
