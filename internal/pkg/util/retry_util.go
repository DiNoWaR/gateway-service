package util

import (
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"io"
	"net/http"
	"time"
)

const callbackHeader = "Callback-URL"

func RetryableRequest(url string, method string, body io.Reader, callbackUrl, contentType string, timeout, waitingTime int) (*http.Response, error) {
	var resp *http.Response
	var err error

	retryTimeout := time.Duration(timeout) * time.Second
	maxElapsedTime := time.Duration(waitingTime) * time.Second

	operation := func() error {
		req, reqErr := http.NewRequest(method, url, body)
		if reqErr != nil {
			return reqErr
		}

		req.Header.Set("Content-Type", contentType)
		req.Header.Set(callbackHeader, callbackUrl)

		client := &http.Client{
			Timeout: retryTimeout,
		}
		resp, err = client.Do(req)

		if err != nil {
			return err
		}

		if resp.StatusCode >= 500 {
			serverErr := fmt.Errorf("server error: %v", resp.Status)
			return serverErr
		}

		return nil
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = maxElapsedTime

	backOffErr := backoff.Retry(operation, expBackoff)
	if backOffErr != nil {
		return nil, backOffErr
	}

	return resp, nil
}
