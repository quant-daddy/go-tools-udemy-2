package tools

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func cleanup(t *testing.T) (string, error) {
	t.Cleanup(func() {
		fmt.Println("Cleanup called")
	})
	return "", nil
}

func Test_ReadJSON(t *testing.T) {
	cleanup(t)
	req, err := http.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte(`{"foo":"bar"}`)))
	if err != nil {
		t.Log("Error", err)
	}
	rr := httptest.NewRecorder()
	var decodedJSON struct {
		Foo string `json:"foo"`
	}
	err = ReadJSON(rr, req, &decodedJSON)
	if err != nil {
		t.Error(err)
	}
	if decodedJSON.Foo != "bar" {
		t.Error("decoding failed")
	}
	fmt.Println("ReadJSON passed")
}

func Test_WriteJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	data := JsonResponse{
		Error:   false,
		Message: "foo",
	}
	headers := make(http.Header)
	headers.Set("FOO", "BAR")
	err := WriteJSON(rr, http.StatusOK, data, headers)
	if err != nil {
		t.Error(err)
	}
	result := JsonResponse{}
	json.Unmarshal(rr.Body.Bytes(), &result)
	if result.Message != "foo" {
		t.Error("data did not match")
	}
}

func Test_WriteError(t *testing.T) {
	rr := httptest.NewRecorder()
	err := WriteError(rr, errors.New("bad"), http.StatusBadGateway)
	if err != nil {
		t.Error(err)
	}
	result := JsonResponse{}
	json.Unmarshal(rr.Body.Bytes(), &result)
	if result.Message != "bad" || !result.Error {
		t.Error("data did not match")
	}
}
