package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// AssertResponse assert response header and body.
func AssertResponse(t *testing.T, res *http.Response, code int, path string) {
	t.Helper()

	AssertResponseHeader(t, res, code)
	AssertResponseBodyWithFile(t, res, path)
}

// AssertResponseHeader assert response header.
func AssertResponseHeader(t *testing.T, res *http.Response, code int) {
	t.Helper()

	// ステータスコードのチェック
	if code != res.StatusCode {
		t.Errorf("expected status code is '%d',\n but actual given code is '%d'", code, res.StatusCode)
	}
	// Content-Typeのチェック
	if expected := "application/json"; res.Header.Get("Content-Type") != expected {
		t.Errorf("unexpected response Content-Type,\n expected: %#v,\n but given #%v", expected, res.Header.Get("Content-Type"))
	}
}

// AssertResponseBodyWithFile assert response body with test file.
func AssertResponseBodyWithFile(t *testing.T, res *http.Response, path string) {
	t.Helper()

	rs := GetStringFromTestFile(t, path)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("unexpected error by ioutil.ReadAll() '%#v'", err)
	}
	var actual bytes.Buffer
	err = json.Indent(&actual, body, "", "  ")
	if err != nil {
		t.Fatalf("unexpected error by json.Indent '%#v'", err)
	}
	assert.JSONEq(t, rs, actual.String())
}

// GetStringFromTestFile get string from test file.
func GetStringFromTestFile(t *testing.T, path string) string {
	t.Helper()

	bt, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("unexpected error while opening file '%#v'", err)
	}
	return string(bt)
}
