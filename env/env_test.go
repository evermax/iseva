package env

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServingJSON(t *testing.T) {
	env := Interface{DB: "db_test.json"}

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

		env.Handler(rec, req)
		if rec.Code != test.expectedStatus {
			t.Fatalf("Test %d, expected status: %d, got %d", i, test.expectedStatus, rec.Code)
		}
		respBody := rec.Body.String()
		if respBody != test.expectedContent {
			t.Fatalf("Test %d, expected body %s, got %s", i, test.expectedContent, respBody)
		}
	}
}
