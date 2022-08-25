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

import "fmt"

// Storage is a wrapper for combined cache and database operations
type Storage struct {
	sqlstorage SQLStorage
	cache      *Cache
}

// Init kicks off the database connector
func (s *Storage) Init(user, password, host, name, redisHost, redisPort string, cache bool) error {
	if err := s.sqlstorage.Init(user, password, host, name); err != nil {
		return err
	}

	var err error
	s.cache, err = NewCache(redisHost, redisPort, cache)
	if err != nil {
		return err
	}

	return nil
}

// List retrieves a list of todos from either cache if cached or database
func (s Storage) List() (Todos, error) {
	ts, err := s.cache.List()
	if err != nil {
		if err == ErrCacheMiss {
			ts, err = s.sqlstorage.List()
			if err != nil {
				return ts, fmt.Errorf("error getting list of todos from database: %v", err)
			}
		}
		if err := s.cache.SaveList(ts); err != nil {
			return ts, fmt.Errorf("error caching list of todos : %v", err)
		}
	}

	return ts, nil
}

// Create records a new todo in the database.
func (s Storage) Create(t Todo) (Todo, error) {
	if err := s.cache.DeleteList(); err != nil {
		return Todo{}, fmt.Errorf("error clearing cache : %v", err)
	}

	t, err := s.sqlstorage.Create(t)
	if err != nil {
		return t, err
	}

	if err = s.cache.Save(t); err != nil {
		return t, err
	}

	return t, nil
}

// Read returns a single todo from cache or database
func (s Storage) Read(id string) (Todo, error) {
	t, err := s.cache.Get(id)
	if err != nil {
		if err == ErrCacheMiss {
			t, err = s.sqlstorage.Read(id)
			if err != nil {
				return t, fmt.Errorf("error getting single from database todo: %v", err)
			}
		}
		if err := s.cache.Save(t); err != nil {
			return t, fmt.Errorf("error caching single todo : %v", err)
		}
	}

	return t, nil
}

// Update changes one todo in the database.
func (s Storage) Update(t Todo) error {
	if err := s.cache.DeleteList(); err != nil {
		return fmt.Errorf("error clearing cache : %v", err)
	}

	if err := s.sqlstorage.Update(t); err != nil {
		return err
	}

	if err := s.cache.Save(t); err != nil {
		return err
	}

	return nil
}

// Delete removes one todo from the database.
func (s Storage) Delete(id string) error {
	if err := s.cache.DeleteList(); err != nil {
		return fmt.Errorf("error clearing cache : %v", err)
	}

	if err := s.sqlstorage.Delete(id); err != nil {
		return err
	}

	if err := s.cache.Delete(id); err != nil {
		return err
	}

	return nil
}
