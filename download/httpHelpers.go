package download

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func doRequest(client *http.Client, req *http.Request, referer, userAgent string) (io.ReadCloser, string, error) {
	req.Header.Set("User-Agent", userAgent)
	if referer != "" {
		req.Header.Set("Referer", referer)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	if resp.StatusCode > 399 {
		resp.Body.Close()
		return nil, "", fmt.Errorf("error: %q %s", resp.Request.URL.String(), resp.Status)
	}
	return resp.Body, resp.Header.Get("Content-Type"), nil
}

func Post(client *http.Client, data url.Values, requestUrl, referer, userAgent string) (io.ReadCloser, string, error) {
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return doRequest(client, req, referer, userAgent)
}

func Get(client *http.Client, requestUrl, referer, userAgent string) (io.ReadCloser, string, error) {
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return nil, "", err
	}
	return doRequest(client, req, referer, userAgent)
}
