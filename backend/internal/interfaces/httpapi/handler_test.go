package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	appcalculation "calculator-app/backend/internal/application/calculation"
)

func request(t *testing.T, method, path, body, contentType string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if contentType != "" { req.Header.Set("Content-Type", contentType) }
	recorder := httptest.NewRecorder()
	NewHandler(appcalculation.NewService()).ServeHTTP(recorder, req)
	return recorder
}

func TestCalculateEndpoint(t *testing.T) {
	response := request(t, "POST", "/api/v1/calculations", `{"operation":"add","a":2,"b":3}`, "application/json")
	if response.Code != http.StatusOK { t.Fatalf("status = %d", response.Code) }
	var body calculationResponse
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil { t.Fatal(err) }
	if body.Result != 5 { t.Errorf("result = %v", body.Result) }
}

func TestCalculateEndpointValidation(t *testing.T) {
	tests := []struct { name, body, contentType string; status int }{
		{"missing content type", `{"operation":"add","a":1,"b":2}`, "", 400}, {"malformed json", `{`, "application/json", 400}, {"empty body", ``, "application/json", 400}, {"two bodies", `{"operation":"add","a":1,"b":2} {}`, "application/json", 400}, {"unknown field", `{"operation":"add","a":1,"b":2,"c":3}`, "application/json", 400}, {"missing a", `{"operation":"sqrt"}`, "application/json", 400}, {"non numeric", `{"operation":"sqrt","a":"x"}`, "application/json", 400}, {"missing b", `{"operation":"add","a":1}`, "application/json", 400}, {"divide by zero", `{"operation":"divide","a":1,"b":0}`, "application/json", 422}, {"bad operation", `{"operation":"wat","a":1}`, "application/json", 400}, {"negative root", `{"operation":"sqrt","a":-1}`, "application/json", 422}, {"zero percentage base", `{"operation":"percentage","a":1,"b":0}`, "application/json", 422}, {"unexpected sqrt operand", `{"operation":"sqrt","a":4,"b":2}`, "application/json", 400},
	}
	for _, test := range tests { t.Run(test.name, func(t *testing.T) { response := request(t, "POST", "/api/v1/calculations", test.body, test.contentType); if response.Code != test.status { t.Errorf("status = %d, want %d; body=%s", response.Code, test.status, response.Body.String()) } }) }
}

func TestHealthAndCORS(t *testing.T) {
	req := httptest.NewRequest("GET", "/healthz", nil)
	response := httptest.NewRecorder()
	NewHandler(appcalculation.NewService()).ServeHTTP(response, req)
	if response.Code != http.StatusOK { t.Errorf("health status = %d", response.Code) }
	if response.Header().Get("Access-Control-Allow-Origin") == "" { t.Error("CORS header missing") }
	options := request(t, "OPTIONS", "/api/v1/calculations", "", "")
	if options.Code != http.StatusNoContent { t.Errorf("OPTIONS status = %d", options.Code) }
}

func TestAPIDocumentation(t *testing.T) {
	for _, path := range []string{"/openapi.yaml", "/docs"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		NewHandler(appcalculation.NewService()).ServeHTTP(response, req)
		if response.Code != http.StatusOK { t.Errorf("%s status = %d", path, response.Code) }
		if response.Body.Len() == 0 { t.Errorf("%s returned an empty body", path) }
	}
}
