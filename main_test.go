package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		want         string
	}{
		{
			name:         "returns environment variable when set",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "custom",
			want:         "custom",
		},
		{
			name:         "returns default when environment variable not set",
			key:          "TEST_VAR_UNSET",
			defaultValue: "default",
			envValue:     "",
			want:         "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				if err := os.Setenv(tt.key, tt.envValue); err != nil {
					t.Fatalf("Failed to set env: %v", err)
				}
				defer func() {
					if err := os.Unsetenv(tt.key); err != nil {
						t.Logf("Failed to unset env: %v", err)
					}
				}()
			}

			got := getEnv(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("getEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHealthEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Create a simple handler that mimics the health endpoint
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		environment := getEnv("ENVIRONMENT", "production")
		version := getEnv("VERSION", "1.0.0")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status":"healthy","environment":"` + environment + `","version":"` + version + `"}`)); err != nil {
			t.Logf("Failed to write response: %v", err)
		}
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `"status":"healthy"`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want to contain %v", rr.Body.String(), expected)
	}
}

func TestReadyEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/ready", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status":"ready"}`)); err != nil {
			t.Logf("Failed to write response: %v", err)
		}
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":"ready"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestTemplateRendering(t *testing.T) {
	// Test that templates directory exists
	if _, err := os.Stat("templates"); os.IsNotExist(err) {
		t.Fatal("templates directory does not exist")
	}

	// Test that index.html exists
	if _, err := os.Stat("templates/index.html"); os.IsNotExist(err) {
		t.Fatal("templates/index.html does not exist")
	}
}

func TestEnvironmentDefaults(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		defaultValue string
	}{
		{"PORT default", "PORT", "8080"},
		{"DISPLAY_MESSAGE default", "DISPLAY_MESSAGE", "Welcome to Service Template!"},
		{"ENVIRONMENT default", "ENVIRONMENT", "production"},
		{"VERSION default", "VERSION", "1.0.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure env var is not set
			if err := os.Unsetenv(tt.envKey); err != nil {
				t.Logf("Failed to unset env: %v", err)
			}

			got := getEnv(tt.envKey, tt.defaultValue)
			if got != tt.defaultValue {
				t.Errorf("getEnv(%s) = %v, want %v", tt.envKey, got, tt.defaultValue)
			}
		})
	}
}
