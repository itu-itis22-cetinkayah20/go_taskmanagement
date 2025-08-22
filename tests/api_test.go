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

	"github.com/gorilla/mux"
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

	// 3) Handler’ları operationId ile eşleyecek registry
	handlerRegistry := map[string]http.HandlerFunc{
		"RegisterHandler":    handlers.RegisterHandler,
		"LoginHandler":       handlers.LoginHandler,
		"PublicTasksHandler": handlers.PublicTasksHandler,
		"TasksListHandler":   handlers.TasksListHandler,
		"TaskCreateHandler":  handlers.TaskCreateHandler,
		"TaskDetailHandler":  handlers.TaskDetailHandler,
		"TaskUpdateHandler":  handlers.TaskUpdateHandler,
		"TaskDeleteHandler":  handlers.TaskDeleteHandler,
		"LogoutHandler":      handlers.LogoutHandler,
	}

	// 4) Router’ı schema’dan dinamik doldur
	router := mux.NewRouter()
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
			// güvenli route’ları sar
			if sec, _ := op["security"].([]interface{}); len(sec) > 0 {
				h = middleware.AuthMiddleware(h)
			}
			router.HandleFunc(path, h).Methods(method)
		}
	}
	for _, p := range static {
		register(p)
	}
	for _, p := range param {
		register(p)
	}

	// 5) Test sunucusunu ayağa kaldır
	srv := httptest.NewServer(router)
	defer srv.Close()

	// 6) Tüm yollar + metod kombinasyonlarını test et
	for path, ops := range spec.Paths {
		for mRaw, rawOp := range ops {
			method := strings.ToUpper(mRaw)
			// alt test
			t.Run(method+" "+path, func(t *testing.T) {
				// a) URL’i hazırla
				url := srv.URL + strings.ReplaceAll(path, "{id}", "1")

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
				// opMap is already unmarshaled above
				if sec, _ := opMap["security"].([]interface{}); len(sec) > 0 {
					rr := httptest.NewRecorder()
					router.ServeHTTP(rr, req)
					if rr.Code != http.StatusUnauthorized {
						t.Errorf("expected 401 for %s %s, got %d", method, path, rr.Code)
					}
					return
				}

				// d) Yetkisiz değilse doğrudan çalıştır ve kodu validate et
				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)

				// spec’deki response kodlarıyla karşılaştır
				expected := map[int]struct{}{}
				if responses, ok := opMap["responses"].(map[string]interface{}); ok {
					for code := range responses {
						if c, err := strconv.Atoi(code); err == nil {
							expected[c] = struct{}{}
						}
					}
				}
				if _, ok := expected[rr.Code]; !ok {
					t.Errorf("%s %s: expected one of %v, got %d",
						method, path, keys(expected), rr.Code)
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
