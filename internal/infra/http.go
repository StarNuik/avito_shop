package infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: this file is ugly

type ErrHandler func(int, []byte) error

type RoundtripHandler func(req *http.Request) (*http.Response, error)

type HttpRequest struct {
	Method         string
	Url            string
	Headers        map[string]string
	In             any
	Out            any
	UnmarshalError ErrHandler
	// TODO: i don't like this
	HttpRoundtrip RoundtripHandler
}

func setupDefaults(r HttpRequest) HttpRequest {
	if len(r.Method) == 0 {
		r.Method = http.MethodGet
	}
	if len(r.Url) == 0 {
		r.Url = "http://localhost"
	}
	if r.UnmarshalError == nil {
		r.UnmarshalError = defaultErrorHandler
	}
	if r.HttpRoundtrip == nil {
		r.HttpRoundtrip = defaultRoundtripHandler
	}
	return r
}

func defaultErrorHandler(code int, _ []byte) error {
	if code/100 != 2 {
		return fmt.Errorf("error")
	}
	return nil
}

func defaultRoundtripHandler(request *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(request)
}

func DoHttp(r HttpRequest) error {
	r = setupDefaults(r)

	reqJson, err := json.Marshal(r.In)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(reqJson)
	request, err := http.NewRequest(r.Method, r.Url, reader)

	for key, value := range r.Headers {
		request.Header.Set(key, value)
	}

	response, err := r.HttpRoundtrip(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if err = r.UnmarshalError(response.StatusCode, body); err != nil {
		return err
	}

	err = json.Unmarshal(body, r.Out)
	if err != nil {
		return err
	}

	return nil
}
