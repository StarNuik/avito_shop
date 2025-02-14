package infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ErrHandler func(int, []byte) error
type HttpHandler func(req *http.Request) (*http.Response, error)

type HttpEngine struct {
	ErrHandler  ErrHandler
	HttpHandler HttpHandler
}

func (e *HttpEngine) Do(method string, url string, headers map[string]string, in any, out any) error {
	reqJson, err := json.Marshal(in)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(reqJson)
	request, err := http.NewRequest(method, url, reader)

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	httpHandler := e.getHttpHandler()
	response, err := httpHandler(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	errHandler := e.getErrHandler()
	if err = errHandler(response.StatusCode, body); err != nil {
		return err
	}

	if out == nil {
		return nil
	}

	err = json.Unmarshal(body, out)
	if err != nil {
		return err
	}

	return nil
}

func (e *HttpEngine) getErrHandler() ErrHandler {
	if e.ErrHandler == nil {
		return defaultErrHandler
	}
	return e.ErrHandler
}

func defaultErrHandler(code int, _ []byte) error {
	if code/100 != 2 {
		return fmt.Errorf("error: %d", code)
	}
	return nil
}

func (e *HttpEngine) getHttpHandler() HttpHandler {
	if e.HttpHandler == nil {
		return defaultHttpHandler
	}
	return e.HttpHandler
}

func defaultHttpHandler(request *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(request)
}
