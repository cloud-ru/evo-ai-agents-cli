package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
)

// Client представляет HTTP клиент для работы с AI Agents API
type Client struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
	projectID  string
}

// NewClient создает новый экземпляр API клиента
func NewClient(baseURL, apiKey, projectID string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey:    apiKey,
		projectID: projectID,
	}
}

// RequestOptions содержит опции для HTTP запроса
type RequestOptions struct {
	Method  string
	Path    string
	Body    interface{}
	Headers map[string]string
	Query   map[string]string
}

// doRequest выполняет HTTP запрос к API
func (c *Client) doRequest(ctx context.Context, opts RequestOptions) (*http.Response, error) {
	url := c.baseURL + opts.Path

	// Добавляем query параметры
	if len(opts.Query) > 0 {
		url += "?"
		first := true
		for key, value := range opts.Query {
			if !first {
				url += "&"
			}
			url += fmt.Sprintf("%s=%s", key, value)
			first = false
		}
	}

	var body io.Reader
	if opts.Body != nil {
		jsonData, err := json.Marshal(opts.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, opts.Method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	for key, value := range opts.Headers {
		req.Header.Set(key, value)
	}

	log.Debug("Making API request", "method", opts.Method, "url", url)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// parseResponse парсит ответ API в указанную структуру
func (c *Client) parseResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errorResp struct {
			Error struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}

		if err := json.Unmarshal(body, &errorResp); err != nil {
			return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
		}

		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, errorResp.Error.Message)
	}

	if target != nil {
		if err := json.Unmarshal(body, target); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}

	return nil
}

// Get выполняет GET запрос
func (c *Client) Get(ctx context.Context, path string, query map[string]string, result interface{}) error {
	resp, err := c.doRequest(ctx, RequestOptions{
		Method: "GET",
		Path:   path,
		Query:  query,
	})
	if err != nil {
		return err
	}

	return c.parseResponse(resp, result)
}

// Post выполняет POST запрос
func (c *Client) Post(ctx context.Context, path string, body interface{}, result interface{}) error {
	resp, err := c.doRequest(ctx, RequestOptions{
		Method: "POST",
		Path:   path,
		Body:   body,
	})
	if err != nil {
		return err
	}

	return c.parseResponse(resp, result)
}

// Put выполняет PUT запрос
func (c *Client) Put(ctx context.Context, path string, body interface{}, result interface{}) error {
	resp, err := c.doRequest(ctx, RequestOptions{
		Method: "PUT",
		Path:   path,
		Body:   body,
	})
	if err != nil {
		return err
	}

	return c.parseResponse(resp, result)
}

// Delete выполняет DELETE запрос
func (c *Client) Delete(ctx context.Context, path string, result interface{}) error {
	resp, err := c.doRequest(ctx, RequestOptions{
		Method: "DELETE",
		Path:   path,
	})
	if err != nil {
		return err
	}

	return c.parseResponse(resp, result)
}
