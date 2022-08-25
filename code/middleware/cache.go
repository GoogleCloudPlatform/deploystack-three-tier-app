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
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/gomodule/redigo/redis"
)

// RedisPool is an interface that allows us to swap in an mock for testing cache
// code.
type RedisPool interface {
	Get() redis.Conn
}

// ErrCacheMiss error indicates that an item is not in the cache
var ErrCacheMiss = fmt.Errorf("item is not in cache")

// NewCache returns an initialized cache ready to go.
func NewCache(redisHost, redisPort string, enabled bool) (*Cache, error) {
	c := &Cache{}
	pool := c.InitPool(redisHost, redisPort)
	c.enabled = enabled
	c.redisPool = pool
	return c, nil
}

// Cache abstracts all of the operations of caching for the application
type Cache struct {
	// redisPool *redis.Pool
	redisPool RedisPool
	enabled   bool
}

func (c *Cache) log(msg string) {
	log.Printf("Cache     : %s\n", msg)
}

// InitPool starts the cache off
func (c Cache) InitPool(redisHost, redisPort string) RedisPool {
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	msg := fmt.Sprintf("Initialized Redis at %s", redisAddr)
	c.log(msg)
	const maxConnections = 10

	pool := redis.NewPool(func() (redis.Conn, error) {
		return redis.Dial("tcp", redisAddr)
	}, maxConnections)

	return pool
}

// Clear removes all items from the cache.
func (c Cache) Clear() error {
	if !c.enabled {
		return nil
	}
	conn := c.redisPool.Get()
	defer conn.Close()

	if _, err := conn.Do("FLUSHALL"); err != nil {
		return err
	}
	return nil
}

// Save records a todo into the cache.
func (c *Cache) Save(todo Todo) error {
	if !c.enabled {
		return nil
	}

	conn := c.redisPool.Get()
	defer conn.Close()

	json, err := todo.JSON()
	if err != nil {
		return fmt.Errorf("cannot convert todo to json: %s", err)
	}

	conn.Send("MULTI")
	conn.Send("SET", strconv.Itoa(todo.ID), json)

	if _, err := conn.Do("EXEC"); err != nil {
		return fmt.Errorf("cannot perform exec operation on cache: %s", err)
	}
	c.log("Successfully saved todo to cache")
	return nil
}

// Get gets a todo from the cache.
func (c *Cache) Get(key string) (Todo, error) {
	t := Todo{}
	if !c.enabled {
		return t, ErrCacheMiss
	}
	conn := c.redisPool.Get()
	defer conn.Close()

	s, err := redis.String(conn.Do("GET", key))
	if err == redis.ErrNil {
		return Todo{}, ErrCacheMiss
	} else if err != nil {
		return Todo{}, err
	}

	if err := json.Unmarshal([]byte(s), &t); err != nil {
		return Todo{}, err
	}
	c.log("Successfully retrieved todo from cache")

	return t, nil
}

// Delete will remove a todo from the cache completely.
func (c *Cache) Delete(key string) error {
	if !c.enabled {
		return nil
	}
	conn := c.redisPool.Get()
	defer conn.Close()

	if _, err := conn.Do("DEL", key); err != nil {
		return err
	}

	c.log(fmt.Sprintf("Cleaning from cache %s", key))
	return nil
}

// List gets all of the todos from the cache.
func (c *Cache) List() (Todos, error) {
	t := Todos{}
	if !c.enabled {
		return t, ErrCacheMiss
	}
	conn := c.redisPool.Get()
	defer conn.Close()

	s, err := redis.String(conn.Do("GET", "todoslist"))
	if err == redis.ErrNil {
		return Todos{}, ErrCacheMiss
	} else if err != nil {
		return Todos{}, err
	}

	if err := json.Unmarshal([]byte(s), &t); err != nil {
		return Todos{}, err
	}
	c.log("Successfully retrieved todos from cache")

	return t, nil
}

// SaveList records a todo list into the cache.
func (c *Cache) SaveList(todos Todos) error {
	if !c.enabled {
		return nil
	}

	conn := c.redisPool.Get()
	defer conn.Close()

	json, err := todos.JSON()
	if err != nil {
		return err
	}

	if _, err := conn.Do("SET", "todoslist", json); err != nil {
		return err
	}
	c.log("Successfully saved todo to cache")
	return nil
}

// DeleteList deletes a todo list into the cache.
func (c *Cache) DeleteList() error {
	if !c.enabled {
		return nil
	}

	return c.Delete("todoslist")
}
