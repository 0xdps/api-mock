package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// NubeUser holds the fields returned by GET /v1/me.
type NubeUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// NubeClient validates session tokens against the nube-auth gateway.
type NubeClient struct {
	gatewayURL string
	httpClient *http.Client
}

// NewNubeClient creates a NubeClient pointing at the given gateway URL
// (e.g. "https://api.nubeauth.com").
func NewNubeClient(gatewayURL string) *NubeClient {
	return &NubeClient{
		gatewayURL: gatewayURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// ValidateSession calls GET /v1/me with the provided Bearer token.
// Returns (nil, nil) when the token is invalid (4xx response).
func (c *NubeClient) ValidateSession(ctx context.Context, token string) (*NubeUser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.gatewayURL+"/v1/me", nil)
	if err != nil {
		return nil, fmt.Errorf("nube: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("nube: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, nil // invalid / expired token — not an error
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("nube: unexpected status %d: %s", resp.StatusCode, body)
	}

	var user NubeUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("nube: decode response: %w", err)
	}
	return &user, nil
}
