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
	"database/sql"
	"strconv"
)

// SQLStorage is a wrapper for database operations
type SQLStorage struct {
	db *sql.DB
}

// Init kicks off the database connector
func (s *SQLStorage) Init(user, password, host, name string) error {
	var err error
	s.db, err = sql.Open("mysql", user+":"+password+"@tcp("+host+")/"+name+"?parseTime=true")
	if err != nil {
		return err
	}

	return nil
}

// Close ends the database connection
func (s *SQLStorage) Close() error {
	return s.db.Close()
}

// List returns a list of all todos
func (s SQLStorage) List() (Todos, error) {
	ts := Todos{}
	results, err := s.db.Query("SELECT * FROM todo ORDER BY updated DESC")
	if err != nil {
		return ts, err
	}

	for results.Next() {
		t, err := resultToTodo(results)
		if err != nil {
			return ts, err
		}

		ts = append(ts, t)
	}
	return ts, nil
}

// Create records a new todo in the database.
func (s SQLStorage) Create(t Todo) (Todo, error) {
	sql := `
		INSERT INTO todo(title, updated) 
		VALUES(?,NOW())	
	`

	if t.Complete {
		sql = `
		INSERT INTO todo(title, updated, completed) 
		VALUES(?,NOW(),NOW())	
	`
	}

	op, err := s.db.Prepare(sql)
	if err != nil {
		return t, err
	}

	results, err := op.Exec(t.Title)

	id, err := results.LastInsertId()
	if err != nil {
		return t, err
	}

	t.ID = int(id)

	return t, nil
}

func resultToTodo(results *sql.Rows) (Todo, error) {
	t := Todo{}
	if err := results.Scan(&t.ID, &t.Title, &t.Updated, &t.completedNull); err != nil {
		return t, err
	}

	if t.completedNull.Valid {
		t.Completed = t.completedNull.Time
		t.Complete = true
	}

	return t, nil
}

// Read returns a single todo from the database
func (s SQLStorage) Read(id string) (Todo, error) {
	t := Todo{}
	results, err := s.db.Query("SELECT * FROM todo WHERE id =?", id)
	if err != nil {
		return t, err
	}

	results.Next()
	t, err = resultToTodo(results)
	if err != nil {
		return t, err
	}

	return t, nil
}

// Update changes one todo in the database.
func (s SQLStorage) Update(t Todo) error {
	orig, err := s.Read(strconv.Itoa(t.ID))
	if err != nil {
		return err
	}

	sql := `
		UPDATE todo
		SET title = ?, updated = NOW() 
		WHERE id = ?
	`

	if t.Complete && !orig.Complete {
		sql = `
		UPDATE todo
		SET title = ?, updated = NOW(), completed = NOW() 
		WHERE id = ?
	`
	}

	if orig.Complete && !t.Complete {
		sql = `
		UPDATE todo
		SET title = ?, updated = NOW(), completed = NULL 
		WHERE id = ?
	`
	}

	op, err := s.db.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = op.Exec(t.Title, t.ID)

	if err != nil {
		return err
	}

	return nil
}

// Delete removes one todo from the database.
func (s SQLStorage) Delete(id string) error {
	op, err := s.db.Prepare("DELETE FROM todo WHERE id =?")
	if err != nil {
		return err
	}

	if _, err = op.Exec(id); err != nil {
		return err
	}

	return nil
}
