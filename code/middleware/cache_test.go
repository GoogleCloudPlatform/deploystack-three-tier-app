package main

import (
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getMockedCache(t *testing.T, enabled bool) (*Cache, *miniredis.Miniredis) {
	s := miniredis.RunT(t)
	sl := strings.Split(s.Addr(), ":")
	c, _ := NewCache(sl[0], sl[1], enabled)

	return c, s
}

func TestCacheCreateReadDelete(t *testing.T) {
	tests := map[string]struct {
		todo Todo
	}{
		"basic": {
			todo: Todo{
				ID:        1,
				Title:     "write a basic test",
				Updated:   time.Date(2023, 03, 19, 10, 30, 0, 0, time.Local),
				Completed: time.Date(2023, 03, 19, 10, 30, 0, 0, time.Local),
				Complete:  true,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			c, s := getMockedCache(t, true)
			_ = s

			err := c.Save(tc.todo)
			require.Nil(t, err)

			got, err := c.Get(tc.todo.Key())
			require.Nil(t, err)

			assert.Equal(t, tc.todo, got)

			err = c.Delete(tc.todo.Key())
			require.Nil(t, err)

			got2, err := c.Get(tc.todo.Key())
			assert.Equal(t, Todo{}, got2)
			assert.Equal(t, ErrCacheMiss, err)

		})
	}
}

func TestCacheCreateReadDeleteCacheError(t *testing.T) {
	tests := map[string]struct {
		todo   Todo
		errStr string
	}{
		"basic": {
			todo: Todo{
				ID:        1,
				Title:     "write a basic test",
				Updated:   time.Date(2023, 03, 19, 10, 30, 0, 0, time.Local),
				Completed: time.Date(2023, 03, 19, 10, 30, 0, 0, time.Local),
				Complete:  true,
			},
			errStr: "connect: connection refused",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			c, s := getMockedCache(t, true)
			s.Close()

			err := c.Save(tc.todo)
			assert.ErrorContains(t, err, tc.errStr)

			_, err = c.Get(tc.todo.Key())
			assert.ErrorContains(t, err, tc.errStr)

			err = c.Delete(tc.todo.Key())
			assert.ErrorContains(t, err, tc.errStr)

			_, err = c.Get(tc.todo.Key())
			assert.ErrorContains(t, err, tc.errStr)

		})
	}
}

func TestCacheCreateReadDeleteDisabled(t *testing.T) {
	tests := map[string]struct {
		todo Todo
	}{
		"basic": {
			todo: Todo{
				ID:        1,
				Title:     "write a basic test",
				Updated:   time.Date(2023, 03, 19, 10, 30, 0, 0, time.Local),
				Completed: time.Date(2023, 03, 19, 10, 30, 0, 0, time.Local),
				Complete:  true,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			c, s := getMockedCache(t, false)
			_ = s

			err := c.Save(tc.todo)
			require.Nil(t, err)

			_, err = c.Get(tc.todo.Key())
			require.Equal(t, ErrCacheMiss, err)

			err = c.Delete(tc.todo.Key())
			require.Nil(t, err)

			_, err = c.Get(tc.todo.Key())
			require.Equal(t, ErrCacheMiss, err)

		})
	}
}

func TestCacheCreateReadDeleteList(t *testing.T) {
	tests := map[string]struct {
		todos Todos
	}{
		"basic": {
			todos: Todos{
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
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			c, s := getMockedCache(t, true)
			_ = s

			err := c.SaveList(tc.todos)
			require.Nil(t, err)

			got, err := c.List()
			require.Nil(t, err)
			assert.Equal(t, tc.todos, got)

			err = c.DeleteList()
			require.Nil(t, err)

			got2, err := c.List()
			assert.Equal(t, Todos{}, got2)
			assert.Equal(t, ErrCacheMiss, err)

		})
	}
}

func TestCacheCreateReadDeleteListError(t *testing.T) {
	tests := map[string]struct {
		todos  Todos
		errStr string
	}{
		"basic": {
			todos: Todos{
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
			},
			errStr: "connect: connection refused",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			c, s := getMockedCache(t, true)
			s.Close()

			err := c.SaveList(tc.todos)
			assert.ErrorContains(t, err, tc.errStr)

			_, err = c.List()
			assert.ErrorContains(t, err, tc.errStr)

			err = c.DeleteList()
			assert.ErrorContains(t, err, tc.errStr)

			_, err = c.List()
			assert.ErrorContains(t, err, tc.errStr)

		})
	}
}

func TestCacheCreateReadDeleteListDisabled(t *testing.T) {
	tests := map[string]struct {
		todos Todos
	}{
		"basic": {
			todos: Todos{
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
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			c, s := getMockedCache(t, false)
			_ = s

			err := c.SaveList(tc.todos)
			assert.Nil(t, err)

			_, err = c.List()
			assert.Equal(t, ErrCacheMiss, err)

			err = c.DeleteList()
			assert.Nil(t, err)

			_, err = c.List()
			assert.Equal(t, ErrCacheMiss, err)

		})
	}
}

func TestCacheClear(t *testing.T) {
	tests := map[string]struct {
		todos Todos
	}{
		"basic": {
			todos: Todos{
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
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			c, s := getMockedCache(t, true)
			_ = s
			c.SaveList(tc.todos)

			got, err := c.List()
			require.Nil(t, err)
			assert.Equal(t, tc.todos, got)

			err = c.Clear()
			require.Nil(t, err)

			got2, err := c.List()
			assert.Equal(t, Todos{}, got2)
			assert.Equal(t, ErrCacheMiss, err)
		})
	}
}

func TestCacheClearError(t *testing.T) {
	tests := map[string]struct {
		todos  Todos
		errStr string
	}{
		"basic": {
			todos: Todos{
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
			},
			errStr: "connect: connection refused",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			c, s := getMockedCache(t, true)
			s.Close()

			err := c.SaveList(tc.todos)
			assert.ErrorContains(t, err, tc.errStr)

			_, err = c.List()
			assert.ErrorContains(t, err, tc.errStr)

			err = c.Clear()
			assert.ErrorContains(t, err, tc.errStr)

			_, err = c.List()
			assert.ErrorContains(t, err, tc.errStr)
		})
	}
}

func TestCacheDisabled(t *testing.T) {
	tests := map[string]struct {
		todos Todos
	}{
		"basic": {
			todos: Todos{
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
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			c, s := getMockedCache(t, false)
			_ = s

			err := c.SaveList(tc.todos)
			assert.Nil(t, err)

			_, err = c.List()
			assert.Equal(t, ErrCacheMiss, err)

			err = c.Clear()
			assert.Nil(t, err)

			_, err = c.List()
			assert.Equal(t, ErrCacheMiss, err)

		})
	}
}
