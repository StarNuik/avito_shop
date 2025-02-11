package infra

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type ErrHandler func(int, []byte) error

func HttpRequest(method string, url string, headers map[string]string, unmarshalError ErrHandler, in any, out any) error {
	reqJson, err := json.Marshal(in)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(reqJson)
	request, err := http.NewRequest(method, url, reader)

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if err = unmarshalError(response.StatusCode, body); err != nil {
		return err
	}

	err = json.Unmarshal(body, out)
	if err != nil {
		return err
	}

	return nil
}
