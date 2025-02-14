package client

import (
	"github.com/avito_shop/internal/infra"
	"github.com/avito_shop/internal/setup"
	"net/http"
	"net/http/httptest"
)

// TODO: move to shoptest (currently not possible because `setup` depends on `shoptest`)
func NewTestClient() Client {
	router := setup.Router()

	return &Impl{
		HostUrl: "",
		HttpEngine: infra.HttpEngine{
			ErrHandler: UnmarshalError,
			HttpHandler: func(req *http.Request) (*http.Response, error) {
				recorder := httptest.NewRecorder()

				router.ServeHTTP(recorder, req)

				return recorder.Result(), nil
			},
		},
	}
}
