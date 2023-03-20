package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errForced = fmt.Errorf("forced error")

type MockCache struct {
	store    map[string]interface{}
	forceErr bool
}

func (m *MockCache) Clear() error {
	if m.forceErr {
		return errForced
	}
	return nil
}
func (m *MockCache) Delete(key string) error {
	if m.forceErr {
		return errForced
	}
	delete(m.store, key)
	return nil
}
func (m *MockCache) DeleteList() error {
	if m.forceErr {
		return errForced
	}
	delete(m.store, "todoslist")
	return nil
}
func (m *MockCache) Get(key string) (Todo, error) {
	if m.forceErr {
		return Todo{}, errForced
	}
	t := m.store[key].(Todo)
	return t, nil
}
func (m *MockCache) InitPool(redisHost string, redisPort string) RedisPool {
	m.store = map[string]interface{}{}
	return &redis.Pool{}
}
func (m *MockCache) List() (Todos, error) {
	if m.forceErr {
		return Todos{}, errForced
	}
	t := m.store["todoslist"].(Todos)
	return t, nil
}
func (m *MockCache) Save(todo Todo) error {
	if m.forceErr {
		return errForced
	}
	m.store[todo.Key()] = todo
	return nil
}
func (m *MockCache) SaveList(todos Todos) error {
	if m.forceErr {
		return errForced
	}
	m.store["todoslist"] = todos
	return nil
}
func (m *MockCache) log(msg string) {
	log.Printf("MockCache : %s\n", msg)
}

type MockDB struct {
	currID   int
	forceErr bool
	todos    map[string]Todo
}

func (m *MockDB) Close() error {
	if m.forceErr {
		return errForced
	}
	return nil
}
func (m *MockDB) Create(t Todo) (Todo, error) {
	if m.forceErr {
		return Todo{}, errForced
	}
	m.currID++
	t.ID = m.currID
	m.todos[t.Key()] = t

	return t, nil
}
func (m *MockDB) Delete(id string) error {
	if m.forceErr {
		return errForced
	}
	delete(m.todos, id)

	return nil
}
func (m *MockDB) Init(user string, password string, host string, name string) error {
	if m.forceErr {
		return errForced
	}
	m.todos = map[string]Todo{}
	return nil
}
func (m *MockDB) List() (Todos, error) {
	if m.forceErr {
		return Todos{}, errForced
	}

	todos := Todos{}
	for _, v := range m.todos {
		todos = append(todos, v)
	}
	return todos, nil
}
func (m *MockDB) Read(id string) (Todo, error) {
	if m.forceErr {
		return Todo{}, errForced
	}
	t := m.todos[id]
	return t, nil
}
func (m *MockDB) Update(t Todo) error {
	if m.forceErr {
		return errForced
	}

	return nil
}

func (s *Storage) MockInit(forceCacheErr, forceSQLErr bool) error {
	cache := MockCache{}
	db := MockDB{}

	cache.InitPool("", "")
	db.Init("", "", "", "")

	cache.forceErr = forceCacheErr
	db.forceErr = forceSQLErr

	s.cache = &cache
	s.sqlstorage = &db
	return nil
}

func TestStorageCreate(t *testing.T) {
	tests := map[string]struct {
		in            Todo
		want          Todo
		forceCacheErr bool
		forceSQLErr   bool
		err           error
	}{
		"basic": {
			in: Todo{
				Title: "write a basic test",
			},
			want: Todo{
				ID:    1,
				Title: "write a basic test",
			},
		},
		"basicCacheErr": {
			in: Todo{
				Title: "write a basic test",
			},
			forceCacheErr: true,
			err:           errForced,
			want:          Todo{},
		},
		"basicSQLErr": {
			in: Todo{
				Title: "write a basic test",
			},
			forceSQLErr: true,
			err:         errForced,
			want:        Todo{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := Storage{}
			s.MockInit(tc.forceCacheErr, tc.forceSQLErr)

			got, err := s.Create(tc.in)
			if tc.err == nil {
				require.Nil(t, err)
			}
			assert.ErrorIs(t, err, tc.err)
			assert.Equal(t, tc.want, got)

		})
	}
}
