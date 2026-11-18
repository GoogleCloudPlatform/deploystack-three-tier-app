// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestParseAllowedOrigins(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want []string
	}{
		{name: "empty", raw: "", want: nil},
		{name: "wildcard rejected", raw: "*", want: nil},
		{name: "single origin with trailing slash", raw: "https://fe.example.run.app/", want: []string{"https://fe.example.run.app"}},
		{name: "multiple origins", raw: "http://localhost:8080, http://127.0.0.1:8080/", want: []string{"http://localhost:8080", "http://127.0.0.1:8080"}},
		{name: "invalid entries ignored", raw: "*, ftp://bad.example, https://good.example", want: []string{"https://good.example"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := parseAllowedOrigins(tc.raw)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("parseAllowedOrigins(%q) = %v, want %v", tc.raw, got, tc.want)
			}
		})
	}
}

func TestRequireAuth(t *testing.T) {
	okHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}

	t.Run("fails closed when expected token is empty", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/todo", nil)
		req.Header.Set("Authorization", "Bearer anything")
		rec := httptest.NewRecorder()

		requireAuth("", okHandler).ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("got status %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("rejects missing header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/todo", nil)
		rec := httptest.NewRecorder()

		requireAuth("secret-token", okHandler).ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("got status %d, want %d", rec.Code, http.StatusUnauthorized)
		}
		if got := rec.Header().Get("WWW-Authenticate"); got == "" {
			t.Fatal("expected WWW-Authenticate header on 401 response")
		}
	})

	t.Run("rejects wrong token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/todo/1", nil)
		req.Header.Set("Authorization", "Bearer wrong-token")
		rec := httptest.NewRecorder()

		requireAuth("secret-token", okHandler).ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("got status %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("allows valid bearer token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/todo", nil)
		req.Header.Set("Authorization", "Bearer secret-token")
		rec := httptest.NewRecorder()

		requireAuth("secret-token", okHandler).ServeHTTP(rec, req)
		if rec.Code != http.StatusCreated {
			t.Fatalf("got status %d, want %d", rec.Code, http.StatusCreated)
		}
	})
}
