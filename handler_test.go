package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServingJSON(t *testing.T) {
	handler := JSONHandler{DB: "testdata/db_simple.json"}
	tests := []struct {
		requestPath     string
		expectedStatus  int
		expectedContent string
	}{
		{
			requestPath:     "/test",
			expectedStatus:  http.StatusOK,
			expectedContent: `{"field1": "value1","field2": "value2"}`,
		},
		{
			requestPath:     "/test/other",
			expectedStatus:  http.StatusOK,
			expectedContent: `{"field3": {"sub3field1": "value1","sub3field2": "value2"},"field4": "value4"}`,
		},
		{
			requestPath:     "/tests",
			expectedStatus:  http.StatusNotFound,
			expectedContent: "",
		},
	}
	for i, test := range tests {
		req, err := http.NewRequest("GET", test.requestPath, nil)
		if err != nil {
			t.Fatalf("An error occured when creating the request: %v for test %d", err, i)
		}
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)
		if rec.Code != test.expectedStatus {
			t.Fatalf("Test %d, expected status: %d, got %d", i, test.expectedStatus, rec.Code)
		}
		respBody := rec.Body.String()
		if respBody != test.expectedContent {
			t.Fatalf("Test %d, expected body %s, got %s", i, test.expectedContent, respBody)
		}
	}
}

func TestErrorGetDBData(t *testing.T) {
	handlers := []JSONHandler{
		{DB: "testdata/db_simple_broken.json"},
		{DB: "testdata/db_second_part_broken.json"},
		{DB: "testdata/db_first_part_broken.json"},
		{DB: "testdata/none.json"},
		{DB: "testdata/db_template_broken.json"},
	}

	for i, handler := range handlers {
		err := handler.getDBData()
		if err == nil {
			t.Fatalf("Test %d: Expected error, got none", i)
		}
	}
}

func TestVariables(t *testing.T) {
	handler := JSONHandler{DB: "testdata/db_variables.json"}
	tests := []struct {
		requestPath     string
		expectedStatus  int
		expectedContent string
	}{
		{
			requestPath:     "/test/variables",
			expectedStatus:  http.StatusOK,
			expectedContent: `{"field5": "test_var1"}`,
		},
	}
	for i, test := range tests {
		req, err := http.NewRequest("GET", test.requestPath, nil)
		if err != nil {
			t.Fatalf("An error occured when creating the request: %v for test %d", err, i)
		}
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)
		if rec.Code != test.expectedStatus {
			t.Fatalf("Test %d, expected status: %d, got %d", i, test.expectedStatus, rec.Code)
		}
		respBody := rec.Body.String()
		if respBody != test.expectedContent {
			t.Fatalf("Test %d, expected body %s, got %s", i, test.expectedContent, respBody)
		}
	}
}

func TestHeaders(t *testing.T) {
	handler := JSONHandler{DB: "testdata/db_simple.json"}
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("An error occured when creating the request: %v", err)
	}
	originHeader := "test.com"
	req.Header.Add("origin", originHeader)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status: %d, got %d", http.StatusOK, rec.Code)
	}
	CORSHeader := rec.Header().Get("Access-Control-Allow-Origin")
	if CORSHeader != originHeader {
		t.Fatalf("Expected allow origin header: %s, got %s", originHeader, CORSHeader)
	}
	expectedAcceptedHeaders := "Content-Type, X-Requested-With"
	respAcceptedHeaders := rec.Header().Get("Access-Control-Allow-Header")
	if expectedAcceptedHeaders != respAcceptedHeaders {
		t.Fatalf("Expected allowed header: %s, got: %s", expectedAcceptedHeaders, respAcceptedHeaders)
	}
	expectedContentType := "application/json; charset=utf-8"
	respContentType := rec.Header().Get("Content-Type")
	if expectedContentType != respContentType {
		t.Fatalf("Expected content type: %s, got: %s", expectedContentType, respContentType)
	}
}

func TestOPTIONSCall(t *testing.T) {
	handler := JSONHandler{DB: "testdata/db_simple.json"}
	req, err := http.NewRequest("OPTIONS", "/test", nil)
	if err != nil {
		t.Fatalf("An error occured when creating the request: %v", err)
	}
	originHeader := "test.com"
	req.Header.Add("origin", originHeader)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("Expected status: %d, got: %d", http.StatusNoContent, rec.Code)
	}
	CORSHeader := rec.Header().Get("Access-Control-Allow-Origin")
	if CORSHeader != originHeader {
		t.Fatalf("Expected allow origin header: %s, got %s", originHeader, CORSHeader)
	}
	expectedAcceptedHeaders := "Content-Type, X-Requested-With"
	respAcceptedHeaders := rec.Header().Get("Access-Control-Allow-Header")
	if expectedAcceptedHeaders != respAcceptedHeaders {
		t.Fatalf("Expected allowed header: %s, got: %s", expectedAcceptedHeaders, respAcceptedHeaders)
	}
	expectedContentType := "application/json; charset=utf-8"
	respContentType := rec.Header().Get("Content-Type")
	if expectedContentType != respContentType {
		t.Fatalf("Expected content type: %s, got: %s", expectedContentType, respContentType)
	}
}

func TestFileNotFound(t *testing.T) {
	handler := JSONHandler{DB: "none.json"}
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("An error occured when creating the request: %v", err)
	}
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("Expected an error status: %d, got: %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestRandom(t *testing.T) {
	handler := JSONHandler{DB: "testdata/db_random.json"}
	req, err := http.NewRequest("GET", "/test/randint", nil)
	if err != nil {
		t.Fatalf("An error occured when creating the request: %v", err)
	}
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected status: %d, got %d", http.StatusOK, rec.Code)
	}
}
