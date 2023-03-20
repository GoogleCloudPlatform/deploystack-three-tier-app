package main

import (
	"fmt"
	"log"
	"testing"
	"time"

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
	tmp := m.store[key]
	if tmp == nil {
		return Todo{}, ErrCacheMiss
	}

	t := tmp.(Todo)
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
	t := Todos{}

	tmp := m.store["todoslist"]

	if tmp == nil {
		return t, ErrCacheMiss
	}

	return tmp.(Todos), nil
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
		log.Printf("%+v", v)
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

func TestStorageList(t *testing.T) {

	todos := Todos{
		Todo{
			ID:        1,
			Title:     "write a basic test",
			Updated:   time.Date(2023, 03, 19, 10, 30, 0, 0, time.Local),
			Completed: time.Date(2023, 03, 19, 10, 30, 0, 0, time.Local),
			Complete:  true,
		},
		Todo{
			ID:        2,
			Title:     "write one more test",
			Updated:   time.Date(2023, 03, 29, 10, 30, 0, 0, time.Local),
			Completed: time.Date(2023, 03, 29, 10, 30, 0, 0, time.Local),
		},
	}

	tests := map[string]struct {
		want          Todos
		forceCacheErr bool
		forceSQLErr   bool
		err           error
	}{
		"basic": {
			want: todos,
		},
		"basicCacheErr": {
			want:          Todos{},
			forceCacheErr: true,
			err:           errForced,
		},
		"basicSQLErr": {
			want:        Todos{},
			forceSQLErr: true,
			err:         errForced,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := Storage{}
			s.MockInit(tc.forceCacheErr, tc.forceSQLErr)

			for _, v := range todos {
				s.Create(v)
			}

			got, err := s.List()
			if tc.err == nil {
				require.Nil(t, err)
			}
			assert.ErrorIs(t, err, tc.err)
			assert.Equal(t, tc.want, got)

		})
	}
}

func TestStorageRead(t *testing.T) {
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
			assert.ErrorIs(t, err, tc.err)
			assert.Equal(t, tc.want, got)

			got2, err := s.Read(got.Key())
			assert.ErrorIs(t, err, tc.err)
			assert.Equal(t, tc.want, got2)

		})
	}
}

func TestStorageDelete(t *testing.T) {
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
			assert.ErrorIs(t, err, tc.err)
			assert.Equal(t, tc.want, got)

			err = s.Delete(got.Key())
			assert.ErrorIs(t, err, tc.err)

			got2, err := s.Read(got.Key())
			assert.ErrorIs(t, err, tc.err)
			assert.Equal(t, Todo{}, got2)

		})
	}
}

func TestStorageUpdate(t *testing.T) {
	tests := map[string]struct {
		in            Todo
		want          Todo
		wantUpdate    Todo
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
			wantUpdate: Todo{
				ID:    1,
				Title: "write another test",
			},
		},
		"basicCacheErr": {
			in: Todo{
				Title: "write a basic test",
			},
			want: Todo{
				ID:    1,
				Title: "write a basic test",
			},
			wantUpdate: Todo{
				ID:    1,
				Title: "write another test",
			},
			forceCacheErr: true,
			err:           errForced,
		},
		"basicSQLErr": {
			in: Todo{
				Title: "write a basic test",
			},
			want: Todo{
				ID:    1,
				Title: "write a basic test",
			},
			wantUpdate: Todo{
				ID:    1,
				Title: "write another test",
			},
			forceSQLErr: true,
			err:         errForced,
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

			if tc.forceCacheErr || tc.forceSQLErr {
				assert.Equal(t, Todo{}, got)
			} else {
				assert.Equal(t, tc.want, got)
			}

			err = s.Update(tc.wantUpdate)
			assert.ErrorIs(t, err, tc.err)

			got2, err := s.Read(got.Key())
			assert.ErrorIs(t, err, tc.err)

			if tc.forceCacheErr || tc.forceSQLErr {
				assert.Equal(t, Todo{}, got2)
			} else {
				assert.Equal(t, tc.wantUpdate, got2)
			}

		})
	}
}
