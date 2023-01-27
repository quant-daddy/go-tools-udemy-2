package tools

import (
	"encoding/json"
	"errors"
	"net/http"
)

type JsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 1 mB
const maxBytes = 1024 * 1024

// ReadJSON reads the json body from request
func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	if data == nil {
		return errors.New("data must not be nil")
	}
	reader := http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(data)
	if err != nil {
		return err
	}
	if decoder.More() {
		return errors.New("only one JSON object allowed in the body")
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for key, val := range headers[0] {
			w.Header()[key] = val
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func WriteError(w http.ResponseWriter, err error, status ...int) error {
	if err == nil {
		return nil
	}
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}
	payload := JsonResponse{
		Error:   true,
		Message: err.Error(),
		Data:    nil,
	}
	return WriteJSON(w, statusCode, payload)
}
