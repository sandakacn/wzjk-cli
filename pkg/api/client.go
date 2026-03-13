package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents the API client
type Client struct {
	baseURL string
	token   string
	client  *http.Client
}

// NewClient creates a new API client
func NewClient(baseURL, token string) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// SetToken updates the authentication token
func (c *Client) SetToken(token string) {
	c.token = token
}

// doRequest performs an HTTP request and decodes the response
func (c *Client) doRequest(method, path string, body interface{}, result interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("error marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp APIResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error != "" {
			return fmt.Errorf("API error: %s", errResp.Error)
		}
		return fmt.Errorf("API error: %s", resp.Status)
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("error parsing response: %w", err)
		}
	}

	return nil
}

// ListDomains returns all domains for the authenticated user
func (c *Client) ListDomains() ([]Domain, error) {
	var resp DomainListResponse
	if err := c.doRequest("GET", "/api/domains", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// AddDomain adds a new domain subscription
func (c *Client) AddDomain(req AddDomainRequest) (*Domain, error) {
	var resp DomainResponse
	if err := c.doRequest("POST", "/api/domains", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// DeleteDomain deletes a domain subscription
func (c *Client) DeleteDomain(id string) error {
	path := fmt.Sprintf("/api/domains?id=%s", id)
	return c.doRequest("DELETE", path, nil, nil)
}

// UpdateDomain updates a domain subscription
func (c *Client) UpdateDomain(id string, req UpdateDomainRequest) (*Domain, error) {
	body := map[string]interface{}{
		"id": id,
	}
	if req.AlertDays != 0 {
		body["alertDays"] = req.AlertDays
	}
	if req.IsActive != nil {
		body["isActive"] = *req.IsActive
	}

	var resp DomainResponse
	if err := c.doRequest("PATCH", "/api/domains", body, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// CheckSSL checks the SSL certificate for a domain
func (c *Client) CheckSSL(domain string, port int) (*SSLInfo, error) {
	req := map[string]interface{}{
		"domain": domain,
	}
	if port > 0 {
		req["port"] = port
	}

	var resp SSLCheckResponse
	if err := c.doRequest("POST", "/api/ssl", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetProfile returns the user profile
func (c *Client) GetProfile() (*UserProfile, error) {
	var resp UserProfileResponse
	if err := c.doRequest("GET", "/api/user/profile", nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetAvailability returns availability status for all domains
func (c *Client) GetAvailability() (AvailabilityResponse, error) {
	var resp struct {
		Success      bool                     `json:"success"`
		Availability AvailabilityResponse     `json:"availability"`
	}
	if err := c.doRequest("GET", "/api/availability", nil, &resp); err != nil {
		return nil, err
	}
	return resp.Availability, nil
}

// LoginResponse represents the response from API key login
type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"user"`
}

// LoginWithAPIKey authenticates using an API key
func (c *Client) LoginWithAPIKey(apiKey string) (*LoginResponse, error) {
	req := map[string]string{
		"apiKey": apiKey,
	}

	var resp struct {
		Success bool          `json:"success"`
		Data    LoginResponse `json:"data"`
		Error   string        `json:"error,omitempty"`
	}

	if err := c.doRequest("POST", "/api/cli/login", req, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}
