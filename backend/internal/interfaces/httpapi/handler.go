package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	appcalculation "calculator-app/backend/internal/application/calculation"
)

type Handler struct {
	service *appcalculation.Service
}

type calculationRequest struct {
	Operation string          `json:"operation"`
	A         json.RawMessage `json:"a"`
	B         json.RawMessage `json:"b"`
}

type calculationResponse struct {
	Operation string  `json:"operation"`
	Result    float64 `json:"result"`
}

type errorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func NewHandler(service *appcalculation.Service) http.Handler {
	handler := &Handler{service: service}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handler.health)
	mux.HandleFunc("POST /api/v1/calculations", handler.calculate)
	return cors(mux)
}

func (handler *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (handler *Handler) calculate(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" && !strings.HasPrefix(contentType, "application/json;") {
		writeError(w, http.StatusBadRequest, "invalid_content_type", "Content-Type must be application/json")
		return
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var request calculationRequest
	if err := decoder.Decode(&request); err != nil {
		message := "request body must be valid JSON"
		if errors.Is(err, io.EOF) { message = "request body is required" }
		writeError(w, http.StatusBadRequest, "invalid_json", message)
		return
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		writeError(w, http.StatusBadRequest, "invalid_json", "request body must contain exactly one JSON object")
		return
	}
	a, err := parseNumber(request.A, "a")
	if err != nil { writeError(w, http.StatusBadRequest, "invalid_operand", err.Error()); return }
	b, err := parseOptionalNumber(request.B, "b")
	if err != nil { writeError(w, http.StatusBadRequest, "invalid_operand", err.Error()); return }
	result, err := handler.service.Calculate(r.Context(), appcalculation.Command{Operation: request.Operation, A: *a, B: b})
	if err != nil {
		status, code := http.StatusUnprocessableEntity, "calculation_error"
		if errors.Is(err, appcalculation.ErrInvalidCommand) {
			status, code = http.StatusBadRequest, "invalid_request"
		}
		writeError(w, status, code, domainMessage(err))
		return
	}
	writeJSON(w, http.StatusOK, calculationResponse{Operation: result.Operation, Result: result.Value})
}

func parseNumber(raw json.RawMessage, name string) (*float64, error) {
	if len(raw) == 0 || string(raw) == "null" { return nil, errors.New(name + " is required") }
	var value float64
	if err := json.Unmarshal(raw, &value); err != nil { return nil, errors.New(name + " must be a number") }
	return &value, nil
}

func parseOptionalNumber(raw json.RawMessage, name string) (*float64, error) {
	if len(raw) == 0 || string(raw) == "null" { return nil, nil }
	return parseNumber(raw, name)
}

func domainMessage(err error) string {
	if wrapped := strings.TrimPrefix(err.Error(), appcalculation.ErrInvalidCommand.Error()+": "); wrapped != err.Error() { return wrapped }
	return err.Error()
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	response := errorResponse{}
	response.Error.Code, response.Error.Message = code, message
	writeJSON(w, status, response)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions { w.WriteHeader(http.StatusNoContent); return }
		next.ServeHTTP(w, r)
	})
}
