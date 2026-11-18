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
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const maxTitleLength = 512

var storage Storage

func main() {
	user := os.Getenv("todo_user")
	pass := os.Getenv("todo_pass")
	host := os.Getenv("todo_host")
	name := os.Getenv("todo_name")
	redisHost := os.Getenv("REDISHOST")
	redisPort := os.Getenv("REDISPORT")
	port := os.Getenv("PORT")
	apiAuthToken := strings.TrimSpace(os.Getenv("API_AUTH_TOKEN"))

	fmt.Printf("Port: %s\n", port)
	if apiAuthToken == "" {
		log.Println("Warning: API_AUTH_TOKEN is not set; state-changing endpoints will reject all requests")
	}

	if err := storage.Init(user, pass, host, name, redisHost, redisPort, true); err != nil {
		panic(err)
	}
	defer storage.sqlstorage.Close()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/healthz", healthHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/healthz", healthHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/todo", listHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/todo", requireAuth(apiAuthToken, createHandler)).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/todo/{id}", readHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/todo/{id}", requireAuth(apiAuthToken, deleteHandler)).Methods(http.MethodDelete)
	router.HandleFunc("/api/v1/todo/{id}", requireAuth(apiAuthToken, updateHandler)).Methods(http.MethodPost, http.MethodPut)

	var handler http.Handler = router
	allowedOrigins := parseAllowedOrigins(os.Getenv("FRONTEND_ORIGIN"))
	if len(allowedOrigins) > 0 {
		headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
		originsOk := handlers.AllowedOrigins(allowedOrigins)
		methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})
		handler = handlers.CORS(originsOk, headersOk, methodsOk)(router)
	} else {
		log.Println("FRONTEND_ORIGIN is not set; cross-origin CORS requests are disabled")
	}

	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func parseAllowedOrigins(raw string) []string {
	var origins []string
	for _, entry := range strings.Split(raw, ",") {
		origin := strings.TrimSpace(entry)
		origin = strings.TrimRight(origin, "/")
		if origin == "" || origin == "*" {
			continue
		}
		u, err := url.Parse(origin)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
			log.Printf("Ignoring invalid FRONTEND_ORIGIN entry: %q", entry)
			continue
		}
		origins = append(origins, u.Scheme+"://"+u.Host)
	}
	return origins
}

func extractBearerToken(r *http.Request) (string, bool) {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader == "" {
		return "", false
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}
	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", false
	}
	return token, true
}

func validAuthToken(provided, expected string) bool {
	if expected == "" || provided == "" {
		return false
	}
	providedHash := sha256.Sum256([]byte(provided))
	expectedHash := sha256.Sum256([]byte(expected))
	return subtle.ConstantTimeCompare(providedHash[:], expectedHash[:]) == 1
}

func requireAuth(expectedToken string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, ok := extractBearerToken(r)
		if !ok || !validAuthToken(token, expectedToken) {
			w.Header().Set("WWW-Authenticate", `Bearer realm="todo-api"`)
			msg := Message{"unauthorized", "valid Bearer token required"}
			writeJSON(w, msg, http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
	return
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	ts, err := storage.List()
	if err != nil {
		writeErrorMsg(w, err)
		return
	}

	writeJSON(w, ts, http.StatusOK)
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	_, err := strconv.Atoi(id)
	if err != nil {
		msg := Message{"invalid! id must be integer", fmt.Sprintf("todo id: %s", id)}
		writeJSON(w, msg, http.StatusInternalServerError)
		return
	}

	t, err := storage.Read(id)
	if err != nil {

		if strings.Contains(err.Error(), "Rows are closed") {
			msg := Message{"todo not found", fmt.Sprintf("todo id: %s", id)}
			writeJSON(w, msg, http.StatusNotFound)
			return
		}

		writeErrorMsg(w, err)
		return
	}

	writeJSON(w, t, http.StatusOK)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	t := Todo{}
	t.Title = r.FormValue("title")
	if utf8.RuneCountInString(t.Title) > maxTitleLength {
		msg := Message{"invalid title", fmt.Sprintf("title must not exceed %d characters", maxTitleLength)}
		writeJSON(w, msg, http.StatusBadRequest)
		return
	}

	if len(r.FormValue("complete")) > 0 && r.FormValue("complete") != "false" {
		t.Complete = true
	}

	t, err := storage.Create(t)
	if err != nil {
		writeErrorMsg(w, err)
		return
	}

	writeJSON(w, t, http.StatusCreated)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	t := Todo{}
	id := mux.Vars(r)["id"]
	t.ID, err = strconv.Atoi(id)
	if err != nil {
		writeErrorMsg(w, err)
		return
	}

	t.Title = r.FormValue("title")
	if utf8.RuneCountInString(t.Title) > maxTitleLength {
		msg := Message{"invalid title", fmt.Sprintf("title must not exceed %d characters", maxTitleLength)}
		writeJSON(w, msg, http.StatusBadRequest)
		return
	}

	if len(r.FormValue("complete")) > 0 && r.FormValue("complete") != "false" {
		t.Complete = true
	}

	if err = storage.Update(t); err != nil {
		writeErrorMsg(w, err)
		return
	}

	writeJSON(w, t, http.StatusOK)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	_, err := strconv.Atoi(id)
	if err != nil {
		msg := Message{"invalid! id must be integer", fmt.Sprintf("todo id: %s", id)}
		writeJSON(w, msg, http.StatusInternalServerError)
		return
	}

	if err := storage.Delete(id); err != nil {
		writeErrorMsg(w, err)
		return
	}
	msg := Message{"todo deleted", fmt.Sprintf("todo id: %s", id)}

	writeJSON(w, msg, http.StatusNoContent)
}

// JSONProducer is an interface that spits out a JSON string version of itself
type JSONProducer interface {
	JSON() (string, error)
	JSONBytes() ([]byte, error)
}

func writeJSON(w http.ResponseWriter, j JSONProducer, status int) {
	json, err := j.JSON()
	if err != nil {
		writeErrorMsg(w, err)
		return
	}
	writeResponse(w, status, json)
	return
}

func writeErrorMsg(w http.ResponseWriter, err error) {
	payload, marshalErr := json.Marshal(map[string]string{"error": err.Error()})
	if marshalErr != nil {
		writeResponse(w, http.StatusInternalServerError, `{"error":"internal server error"}`)
		return
	}
	writeResponse(w, http.StatusInternalServerError, string(payload))
}

func writeResponse(w http.ResponseWriter, status int, msg string) {
	if status != http.StatusOK {
		weblog(msg)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(msg))

	return
}

func weblog(msg string) {
	log.Printf("Webserver : %s", msg)
}

// Message is a structure for communicating additional data to API consumer.
type Message struct {
	Text    string `json:"text"`
	Details string `json:"details"`
}

// JSON marshalls the content of a todo to json.
func (m Message) JSON() (string, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("could not marshal json for response: %s", err)
	}

	return string(bytes), nil
}

// JSONBytes marshalls the content of a todo to json as a byte array.
func (m Message) JSONBytes() ([]byte, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return []byte{}, fmt.Errorf("could not marshal json for response: %s", err)
	}

	return bytes, nil
}
