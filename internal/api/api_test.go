package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/miriyaladhanwinn/on-duty/internal/db"
)

func TestHealthEndpoint(t *testing.T) {
	store := db.NewMemoryStore()
	server := NewServer(store)

	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got: %d", w.Code)
	}
}

func TestLoginAndGetCategories(t *testing.T) {
	store := db.NewMemoryStore()
	server := NewServer(store)

	// 1. Login with seeded student credentials
	loginPayload := map[string]string{
		"email":    "dhanwinn@cb.students.amrita.edu",
		"password": "Password@123",
	}
	body, _ := json.Marshal(loginPayload)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Login failed: status %d, body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	token, ok := resp["token"].(string)
	if !ok || token == "" {
		t.Fatalf("Login response did not return token: %v", resp)
	}

	// 2. Fetch categories
	catReq, _ := http.NewRequest("GET", "/api/v1/categories", nil)
	catW := httptest.NewRecorder()
	server.Router().ServeHTTP(catW, catReq)

	if catW.Code != http.StatusOK {
		t.Fatalf("Categories endpoint failed: status %d", catW.Code)
	}

	// 3. Check /me with Bearer token
	meReq, _ := http.NewRequest("GET", "/api/v1/auth/me", nil)
	meReq.Header.Set("Authorization", "Bearer "+token)
	meW := httptest.NewRecorder()
	server.Router().ServeHTTP(meW, meReq)

	if meW.Code != http.StatusOK {
		t.Fatalf("/me failed: status %d, body: %s", meW.Code, meW.Body.String())
	}
}
