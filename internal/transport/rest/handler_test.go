package rest

import (
	"EcomTechGo/internal/models"
	"EcomTechGo/internal/store"
	"bytes"
	"context"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_Create(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		expectedCode int
	}{
		{
			name:         "Success",
			body:         `{"title": "Buy milk", "description": "urgent", "completed": false}`,
			expectedCode: 201,
		},
		{
			name:         "Empty Title",
			body:         `{"title": "", "description": "oops", "completed": false}`,
			expectedCode: 400,
		},
		{
			name:         "Invalid JSON",
			body:         `{broken_json`,
			expectedCode: 400,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			st := store.New()
			h := NewHandler(st)

			req := httptest.NewRequest("POST", "/todos", bytes.NewBufferString(tc.body))
			rr := httptest.NewRecorder()

			h.ServeHTTP(rr, req)

			if rr.Code != tc.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tc.expectedCode)
			}
		})
	}
}

func TestHandler_GetOne(t *testing.T) {
	st := store.New()
	if err := st.Create(context.Background(), &models.Todo{Title: "Existing Task", ID: 1}); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	h := NewHandler(st)

	tests := []struct {
		name         string
		url          string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Get Existing ID 1",
			url:          "/todos/1",
			expectedCode: 200,
			expectedBody: "Existing Task",
		},
		{
			name:         "Get Non-Existing ID 999",
			url:          "/todos/999",
			expectedCode: 404,
			expectedBody: "Not Found",
		},
		{
			name:         "Invalid ID Format",
			url:          "/todos/abc",
			expectedCode: 400,
			expectedBody: "Invalid ID",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.url, nil)
			rr := httptest.NewRecorder()

			h.ServeHTTP(rr, req)

			if rr.Code != tc.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tc.expectedCode)
			}

			if !strings.Contains(rr.Body.String(), tc.expectedBody) {
				t.Errorf("handler returned unexpected body: got %v want substring %v",
					rr.Body.String(), tc.expectedBody)
			}
		})
	}
}

func TestHandler_Update(t *testing.T) {
	st := store.New()
	if err := st.Create(context.Background(), &models.Todo{Title: "Old Title", ID: 1}); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	h := NewHandler(st)

	tests := []struct {
		name         string
		url          string
		body         string
		expectedCode int
	}{
		{
			name:         "Update Success",
			url:          "/todos/1",
			body:         `{"title": "New Title", "completed": true}`,
			expectedCode: 200,
		},
		{
			name:         "Update Non-Existing",
			url:          "/todos/999",
			body:         `{"title": "Ghost", "completed": true}`,
			expectedCode: 404,
		},
		{
			name:         "Update Invalid Data (Empty Title)",
			url:          "/todos/1",
			body:         `{"title": "", "completed": true}`,
			expectedCode: 400,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("PUT", tc.url, bytes.NewBufferString(tc.body))
			rr := httptest.NewRecorder()

			h.ServeHTTP(rr, req)

			if rr.Code != tc.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tc.expectedCode)
			}
		})
	}
}

func TestHandler_Delete(t *testing.T) {
	st := store.New()
	if err := st.Create(context.Background(), &models.Todo{Title: "To be deleted", ID: 1}); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	h := NewHandler(st)

	tests := []struct {
		name         string
		url          string
		expectedCode int
	}{
		{
			name:         "Delete Success",
			url:          "/todos/1",
			expectedCode: 204,
		},
		{
			name:         "Delete Non-Existing",
			url:          "/todos/999",
			expectedCode: 404,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", tc.url, nil)
			rr := httptest.NewRecorder()

			h.ServeHTTP(rr, req)

			if rr.Code != tc.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					rr.Code, tc.expectedCode)
			}
		})
	}
}
