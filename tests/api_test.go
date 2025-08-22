// // filepath: tests/api_test.go
package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"

	"go_taskmanagement/handlers"
	"go_taskmanagement/middleware"
	"go_taskmanagement/models"

	"github.com/gofiber/fiber/v2"
)

func TestSchemaDrivenAPI(t *testing.T) {
	// 1) Şemayı oku
	data, err := os.ReadFile("../docs/swagger.json")
	if err != nil {
		t.Fatalf("failed to read swagger.json: %v", err)
	}
	var spec struct {
		Paths map[string]map[string]json.RawMessage `json:"paths"`
	}
	if err := json.Unmarshal(data, &spec); err != nil {
		t.Fatalf("invalid swagger JSON: %v", err)
	}

	// 2) Uygulama durumunu sıfırla
	models.Users = []models.User{}
	models.Tasks = []models.Task{}

	// 3) otomatik handler registry (operationId eşlemesi)
	handlerRegistry := handlers.OperationRegistry

	// 4) Fiber app oluşturup schema’dan dinamik doldur
	app := fiber.New()
	static, param := []string{}, []string{}
	for p := range spec.Paths {
		if strings.Contains(p, "{") {
			param = append(param, p)
		} else {
			static = append(static, p)
		}
	}
	sort.Strings(static)
	sort.Strings(param)
	register := func(path string) {
		for mRaw, rawOp := range spec.Paths[path] {
			method := strings.ToUpper(mRaw)
			var op map[string]interface{}
			json.Unmarshal(rawOp, &op)
			opID, _ := op["operationId"].(string)
			h, ok := handlerRegistry[opID]
			if !ok {
				continue
			}
			// fiber path (convert {id} to :id)
			fp := strings.ReplaceAll(strings.ReplaceAll(path, "{", ":"), "}", "")
			// güvenli route’ları sar
			if sec, _ := op["security"].([]interface{}); len(sec) > 0 {
				app.Add(method, fp, middleware.AuthMiddleware, h)
			} else {
				app.Add(method, fp, h)
			}
		}
	}
	for _, p := range static {
		register(p)
	}
	for _, p := range param {
		register(p)
	}

	// test için Fiber app hazır

	// 6) Tüm yollar + metod kombinasyonlarını test et
	for path, ops := range spec.Paths {
		for mRaw, rawOp := range ops {
			method := strings.ToUpper(mRaw)
			// alt test
			t.Run(method+" "+path, func(t *testing.T) {
				// a) URL’i hazırla
				url := strings.ReplaceAll(path, "{id}", "1")

				// b) Body gerekirse {} at
				var body io.Reader
				var opMap map[string]interface{}
				json.Unmarshal(rawOp, &opMap)
				if _, hasBody := opMap["requestBody"]; hasBody {
					// eğer şema body tanımı varsa (requestBody)
					body = strings.NewReader(`{}`)
				}

				req := httptest.NewRequest(method, url, body)
				if body != nil {
					req.Header.Set("Content-Type", "application/json")
				}
				// isteği Fiber app üzerinde çalıştır
				resp, err := app.Test(req, 1000)
				if err != nil {
					t.Fatalf("request failed: %v", err)
				}
				// unauthorized kontrolü
				if sec, _ := opMap["security"].([]interface{}); len(sec) > 0 {
					if resp.StatusCode != http.StatusUnauthorized {
						t.Errorf("expected 401 for %s %s, got %d", method, path, resp.StatusCode)
					}
					return
				}

				// spec’deki response kodlarıyla karşılaştır
				expected := map[int]struct{}{}
				if responses, ok := opMap["responses"].(map[string]interface{}); ok {
					for code := range responses {
						if c, err := strconv.Atoi(code); err == nil {
							expected[c] = struct{}{}
						}
					}
				}
				if _, ok := expected[resp.StatusCode]; !ok {
					t.Errorf("%s %s: expected one of %v, got %d",
						method, path, keys(expected), resp.StatusCode)
				}
			})
		}
	}
}

// keys map[int]struct{} → []int helper
func keys(m map[int]struct{}) []int {
	out := make([]int, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Ints(out)
	return out
}
