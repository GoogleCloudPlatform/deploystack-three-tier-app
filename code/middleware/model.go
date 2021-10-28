package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Todo struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Updated       time.Time `json:"updated"`
	Completed     time.Time `json:"completed"`
	Complete      bool      `json:"complete"`
	completedNull sql.NullTime
}

// JSON marshalls the content of a todo to json.
func (t Todo) JSON() (string, error) {
	bytes, err := json.Marshal(t)
	if err != nil {
		return "", fmt.Errorf("could not marshal json for response: %s", err)
	}

	return string(bytes), nil
}

// JSONBytes marshalls the content of a todo to json as a byte array.
func (t Todo) JSONBytes() ([]byte, error) {
	bytes, err := json.Marshal(t)
	if err != nil {
		return []byte{}, fmt.Errorf("could not marshal json for response: %s", err)
	}

	return bytes, nil
}

// Key returns the id as a string.
func (t Todo) Key() string {
	return strconv.Itoa(t.ID)
}

type Todos []Todo

// JSON marshalls the content of a slice of todos to json.
func (t Todos) JSON() (string, error) {
	bytes, err := json.Marshal(t)
	if err != nil {
		return "", fmt.Errorf("could not marshal json for response: %s", err)
	}

	return string(bytes), nil
}

// JSONBytes marshalls the content of a slice of todos to json as a byte array.
func (t Todos) JSONBytes() ([]byte, error) {
	bytes, err := json.Marshal(t)
	if err != nil {
		return []byte{}, fmt.Errorf("could not marshal json for response: %s", err)
	}

	return bytes, nil
}
