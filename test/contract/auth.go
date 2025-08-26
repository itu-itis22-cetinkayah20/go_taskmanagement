package contract

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type TokenSource interface {
	// Returns bearer token (without "Bearer " prefix). Cached.
	BearerToken() (string, error)
}

type StaticToken string

func (s StaticToken) BearerToken() (string, error) { return string(s), nil }

// JWT bootstrap via your login endpoint
type LoginToken struct {
	BaseURL string
	Email   string // Changed from User to Email to match your API
	Pass    string
	App     *fiber.App // For in-process testing

	mu    sync.Mutex
	token string
	exp   time.Time
}

func (l *LoginToken) BearerToken() (string, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if time.Now().Before(l.exp) && l.token != "" {
		return l.token, nil
	}

	// Match your login request format
	body, _ := json.Marshal(map[string]string{
		"email":    l.Email, // Your API uses email, not username
		"password": l.Pass,
	})

	// Create request for Fiber in-process testing
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	var resp *http.Response
	var err error

	if l.App != nil {
		// Use Fiber in-process testing
		resp, err = l.App.Test(req, 10_000)
	} else {
		// Use HTTP client (for external testing)
		req, _ = http.NewRequest(http.MethodPost, l.BaseURL+"/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err = http.DefaultClient.Do(req)
	}

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Match your LoginResponse format from OpenAPI spec
	var out struct {
		Token string `json:"token"` // Your API returns "token", not "access_token"
		User  struct {
			ID       int    `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
		} `json:"user"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}

	l.token = out.Token
	// Set expiration (JWT tokens typically expire in 1 hour, adjust as needed)
	l.exp = time.Now().Add(55 * time.Minute) // 55 minutes to account for clock skew

	return l.token, nil
}

// Helper function to create a login token for testing
func NewLoginToken(baseURL, email, password string) *LoginToken {
	return &LoginToken{
		BaseURL: baseURL,
		Email:   email,
		Pass:    password,
	}
}

// Helper function for quick static token creation
func NewStaticToken(token string) TokenSource {
	return StaticToken(token)
}
