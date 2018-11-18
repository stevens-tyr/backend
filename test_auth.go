package main

import (
  "bytes"
  "encoding/json"
  "net/http"
  "net/http/httptest"
  "testing"
)

func performRequest(r http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
  var req *http.Request
  if method == "GET" {
    req, _ = http.NewRequest(method, path, nil)
  } else {
    bb := new(bytes.Buffer)
    json.NewEncoder(bb).Encode(body)
    req, _ = http.NewRequest(method, path, bytes.NewBuffer(body))
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Content-Type", "application/json; charset=utf-8")
  }
  resp := httptest.NewRecorder()
  r.ServeHTTP(resp, req)
  return resp
}

func TestInvalidLogin(t *testing.T) {

}
