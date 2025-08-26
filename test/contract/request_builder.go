package contract

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
)

type BuildInput struct {
	Doc              *openapi3.T
	PathTmpl         string
	Op               *openapi3.Operation
	Method           string
	BaseURL          string // optional when using Fiber in-process
	Token            TokenSource
	LoginCredentials *LoginCredentials // For login endpoint testing
}

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Returns *http.Request ready to send (with path, query, headers, and body populated).
func BuildRequest(ctx context.Context, in BuildInput) (*http.Request, error) {
	path := in.PathTmpl

	// 1) Path params
	for _, p := range in.Op.Parameters {
		if p.Value == nil || p.Value.In != openapi3.ParameterInPath {
			continue
		}
		name := p.Value.Name
		val := exampleForSchema(p.Value.Schema)
		path = strings.ReplaceAll(path, "{"+name+"}", url.PathEscape(fmt.Sprintf("%v", val)))
	}

	// 2) Query params
	q := url.Values{}
	for _, p := range in.Op.Parameters {
		if p.Value == nil || p.Value.In != openapi3.ParameterInQuery {
			continue
		}
		name := p.Value.Name
		val := exampleForSchema(p.Value.Schema)
		if val != nil {
			q.Set(name, fmt.Sprintf("%v", val))
		}
	}

	// 3) Headers
	h := http.Header{}
	for _, p := range in.Op.Parameters {
		if p.Value == nil || p.Value.In != openapi3.ParameterInHeader {
			continue
		}
		name := p.Value.Name
		val := exampleForSchema(p.Value.Schema)
		if val != nil {
			h.Set(name, fmt.Sprintf("%v", val))
		}
	}

	// 4) Body (JSON only for brevity; extend as needed)
	var bodyBytes []byte
	if in.Op.RequestBody != nil && in.Op.RequestBody.Value != nil {
		for mt, mtObj := range in.Op.RequestBody.Value.Content {
			if strings.HasPrefix(mt, "application/json") && mtObj.Schema != nil {
				var doc map[string]any

				// Special case for login endpoint
				if in.Method == "POST" && in.PathTmpl == "/login" && in.LoginCredentials != nil {
					doc = map[string]any{
						"email":    in.LoginCredentials.Email,
						"password": in.LoginCredentials.Password,
					}
				} else {
					doc = jsonForSchema(mtObj.Schema)
				}

				bodyBytes, _ = json.Marshal(doc)
				h.Set("Content-Type", "application/json")
				break
			}
		}
	}

	// 5) Security – add Bearer if required
	if requiresAuth(in.Doc, in.Op) && in.Token != nil {
		tok, err := in.Token.BearerToken()
		if err == nil && tok != "" {
			h.Set("Authorization", "Bearer "+tok)
		}
	}

	u := path
	if in.BaseURL != "" {
		u = strings.TrimRight(in.BaseURL, "/") + path
	}
	if len(q) > 0 {
		if strings.Contains(u, "?") {
			u += "&" + q.Encode()
		} else {
			u += "?" + q.Encode()
		}
	}

	req, _ := http.NewRequest(strings.ToUpper(in.Method), u, bytes.NewReader(bodyBytes))
	req = req.WithContext(ctx)
	req.Header = h
	return req, nil
}

func requiresAuth(doc *openapi3.T, op *openapi3.Operation) bool {
	// Operation-level security overrides/global.
	if op.Security != nil && len(*op.Security) > 0 {
		return true
	}
	return len(doc.Security) > 0
}

func getSchemaType(s *openapi3.Schema) string {
	if s.Type != nil && len(*s.Type) > 0 {
		return (*s.Type)[0] // Get first type from the slice
	}
	return "string" // Default fallback
}

func exampleForSchema(ref *openapi3.SchemaRef) any {
	if ref == nil || ref.Value == nil {
		return "1"
	}
	s := ref.Value
	if s.Example != nil {
		return s.Example
	}
	if len(s.Enum) > 0 {
		return s.Enum[0]
	}
	if s.Default != nil {
		return s.Default
	}

	schemaType := getSchemaType(s)
	switch schemaType {
	case "string":
		if s.Format == "uuid" {
			return "00000000-0000-0000-0000-000000000000"
		}
		if s.Format == "date-time" {
			return "2024-01-01T00:00:00Z"
		}
		if s.Format == "email" {
			// Generate unique email for testing to avoid conflicts
			return fmt.Sprintf("test_%s@example.com", generateRandomString())
		}
		// Generate unique string values for potential username fields
		return fmt.Sprintf("user_%s", generateRandomString())
	case "integer":
		return 1
	case "number":
		return 1.0
	case "boolean":
		return true
	case "array":
		return []any{exampleForSchema(s.Items)}
	case "object":
		out := map[string]any{}
		for k, v := range s.Properties {
			out[k] = exampleForSchema(v)
		}
		return out
	default:
		return "example"
	}
}

func jsonForSchema(ref *openapi3.SchemaRef) map[string]any {
	if ref == nil || ref.Value == nil {
		return nil
	}
	s := ref.Value
	schemaType := getSchemaType(s)

	if schemaType == "object" {
		m := map[string]any{}
		// required first
		for _, rk := range s.Required {
			if prop, ok := s.Properties[rk]; ok {
				m[rk] = exampleForSchema(prop)
			}
		}
		// then optional props (kept minimal)
		for k, v := range s.Properties {
			if _, already := m[k]; already {
				continue
			}
			// keep small to avoid heavy payloads
			if v.Value != nil && v.Value.ReadOnly {
				continue
			}
			m[k] = exampleForSchema(v)
		}
		return m
	}
	// For non-objects, wrap under a generic key when body expects primitives/arrays (rare).
	val := exampleForSchema(ref)
	switch reflect.ValueOf(val).Kind() {
	case reflect.Map:
		return val.(map[string]any)
	default:
		return map[string]any{"value": val}
	}
}

// generateRandomString creates a random string for unique test data
func generateRandomString() string {
	// Use crypto/rand with larger byte array and add nanosecond timestamp
	bytes := make([]byte, 16) // Increased to 16 bytes for more uniqueness
	rand.Read(bytes)
	// Add nanosecond timestamp to ensure uniqueness even if rand produces same value
	nanos := fmt.Sprintf("%d", time.Now().UnixNano())
	result := hex.EncodeToString(bytes) + nanos[len(nanos)-8:] // Take last 8 digits
	return result
}
