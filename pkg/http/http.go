package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// HTTPClient is a struct for making HTTP requests
type HTTPClient struct {
	baseURL  string
	method   string
	path     string
	formData bool // Flag to indicate whether form data is required
	headers  map[string]string
	body     map[string]interface{}
	client   *http.Client
}

// NewHTTPClient creates a new HTTPClient instance
func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		baseURL:  baseURL,
		client:   &http.Client{},
		headers:  make(map[string]string),
		formData: false, // Default to JSON data
	}
}

// Request sends an HTTP request with the specified method, path, body, and headers
func (c *HTTPClient) Request(ctx context.Context) ([]byte, error) {
	var (
		req *http.Request
		err error
	)

	urlString := fmt.Sprintf("%s%s", c.baseURL, c.path)

	// Handle form data if required
	if c.formData {
		formData := url.Values{}
		for key, value := range c.body {
			formData.Set(key, value.(string))
		}

		req, err = http.NewRequestWithContext(ctx, c.method, urlString, strings.NewReader(formData.Encode()))
		if err != nil {
			return nil, err
		}

		// Set Content-Type header for form data
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		// Encode JSON data
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(c.body); err != nil {
			return nil, err
		}

		req, err = http.NewRequestWithContext(ctx, c.method, urlString, &buf)
		if err != nil {
			return nil, err
		}

		// Set Content-Type header for JSON data
		req.Header.Set("Content-Type", "application/json")
	}

	// Set request headers
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func (c *HTTPClient) Method(method string) *HTTPClient {
	c.method = method
	return c
}

func (c *HTTPClient) Body(body map[string]interface{}) *HTTPClient {
	c.body = body
	return c
}

func (c *HTTPClient) Path(path string) *HTTPClient {
	c.path = path
	return c
}

func (c *HTTPClient) Headers(headers map[string]string) *HTTPClient {
	c.headers = headers
	return c
}

// SetFormData sets the flag indicating whether form data is required
func (c *HTTPClient) SetFormData(formData bool) *HTTPClient {
	c.formData = formData
	return c
}
