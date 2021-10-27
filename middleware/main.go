package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var storage Storage

func main() {
	user := os.Getenv("todo_user")
	pass := os.Getenv("todo_pass")
	host := os.Getenv("todo_host")
	name := os.Getenv("todo_name")
	port := os.Getenv("PORT")

	fmt.Printf("Port: %s\n", port)

	if err := storage.Init(user, pass, host, name); err != nil {
		panic(err)
	}
	defer storage.Close()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/v1/todo", listHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/todo", createHandler).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/todo/{id}", readHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/todo/{id}", deleteHandler).Methods(http.MethodDelete)
	router.HandleFunc("/api/v1/todo/{id}", updateHandler).Methods(http.MethodPost, http.MethodPut)
	log.Fatal(http.ListenAndServe(":"+port, router))
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
	t.Description = r.FormValue("description")

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
	t.Description = r.FormValue("description")

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
	writeResponse(w, http.StatusOK, json)
	return
}

func writeErrorMsg(w http.ResponseWriter, err error) {
	s := fmt.Sprintf("{\"error\":\"%s\"}", err)
	writeResponse(w, http.StatusInternalServerError, s)
	return
}

func writeResponse(w http.ResponseWriter, status int, msg string) {
	if status != http.StatusOK {
		weblog(fmt.Sprintf(msg))
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("Access-Control-Allow-Origin", "*")
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
