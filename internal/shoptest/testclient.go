package shoptest

import (
	"github.com/avito_shop/internal/client"
	"github.com/avito_shop/internal/infra"
	"github.com/avito_shop/internal/setup"
	"net/http"
	"net/http/httptest"
)

func NewTestClient() client.Client {
	router := setup.Router()

	return &client.Impl{
		HostUrl: "",
		HttpEngine: infra.HttpEngine{
			ErrHandler: client.UnmarshalError,
			HttpHandler: func(req *http.Request) (*http.Response, error) {
				recorder := httptest.NewRecorder()

				router.ServeHTTP(recorder, req)

				return recorder.Result(), nil
			},
		},
	}
}
