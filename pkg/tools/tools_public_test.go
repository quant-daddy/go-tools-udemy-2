package tools_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	tools "github.com/quant-daddy/go-tools-udemy-2/pkg/tools"
)

func Test_ReadJSON(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte(`{"foo":"bar"}`)))
	if err != nil {
		t.Log("Error", err)
	}
	rr := httptest.NewRecorder()
	var decodedJSON struct {
		Foo string `json:"foo"`
	}
	err = tools.ReadJSON(rr, req, &decodedJSON)
	if err != nil {
		t.Error(err)
	}
	if decodedJSON.Foo != "bar" {
		t.Error("decoding failed")
	}
	fmt.Println("ReadJSON passed")
}
