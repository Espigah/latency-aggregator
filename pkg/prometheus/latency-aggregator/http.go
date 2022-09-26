package latencyaggregator

import (
	"encoding/json"
	"io"
	"net/http"
)

func request(method, url string, body io.Reader, data interface{}) (*http.Response, error) {

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
