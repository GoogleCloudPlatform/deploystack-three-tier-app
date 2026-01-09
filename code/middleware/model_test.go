package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTodoJson(t *testing.T) {
	tests := map[string]struct {
		in   Todo
		want string
	}{
		"basic": {
			in: Todo{
				ID:        1,
				Title:     "write a basic test",
				Updated:   time.Date(2023, 03, 19, 10, 30, 0, 0, time.Local),
				Completed: time.Date(2023, 03, 19, 10, 30, 0, 0, time.Local),
				Complete:  true,
			},
			want: `{"id":1,"title":"write a basic test","updated":"2023-03-19T10:30:00-07:00","completed":"2023-03-19T10:30:00-07:00","complete":true}`,
		},
		"empty": {
			in:   Todo{},
			want: `{"id":0,"title":"","updated":"0001-01-01T00:00:00Z","completed":"0001-01-01T00:00:00Z","complete":false}`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.in.JSON()

			require.Nil(t, err)
			assert.Equal(t, tc.want, got)

		})
	}
}

func TestTodoJsonBytes(t *testing.T) {
	tests := map[string]struct {
		in   Todo
		want string
	}{
		"basic": {
			in: Todo{
				ID:        1,
				Title:     "write a basic test",
				Updated:   time.Date(2023, 03, 19, 10, 30, 0, 0, time.Local),
				Completed: time.Date(2023, 03, 19, 10, 30, 0, 0, time.Local),
				Complete:  true,
			},
			want: `{"id":1,"title":"write a basic test","updated":"2023-03-19T10:30:00-07:00","completed":"2023-03-19T10:30:00-07:00","complete":true}`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.in.JSONBytes()

			require.Nil(t, err)
			assert.Equal(t, []byte(tc.want), got)

		})
	}
}

func TestTodoKey(t *testing.T) {
	tests := map[string]struct {
		in   Todo
		want string
	}{
		"basic": {
			in: Todo{
				ID:        1,
				Title:     "write a basic test",
				Updated:   time.Date(2023, 03, 19, 10, 30, 0, 0, time.Local),
				Completed: time.Date(2023, 03, 19, 10, 30, 0, 0, time.Local),
				Complete:  true,
			},
			want: "1",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.in.Key()
			assert.Equal(t, tc.want, got)

		})
	}
}

func TestTodosJson(t *testing.T) {
	tests := map[string]struct {
		in   Todos
		want string
	}{
		"basic": {
			in: Todos{
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

			want: `[{"id":1,"title":"write a basic test","updated":"2023-03-19T10:30:00-07:00","completed":"2023-03-19T10:30:00-07:00","complete":true},{"id":2,"title":"write one more test","updated":"2023-03-29T10:30:00-07:00","completed":"2023-03-29T10:30:00-07:00","complete":false}]`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.in.JSON()

			require.Nil(t, err)
			assert.Equal(t, tc.want, got)

		})
	}
}

func TestTodosJsonBytes(t *testing.T) {
	tests := map[string]struct {
		in   Todos
		want string
	}{
		"basic": {
			in: Todos{
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

			want: `[{"id":1,"title":"write a basic test","updated":"2023-03-19T10:30:00-07:00","completed":"2023-03-19T10:30:00-07:00","complete":true},{"id":2,"title":"write one more test","updated":"2023-03-29T10:30:00-07:00","completed":"2023-03-29T10:30:00-07:00","complete":false}]`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.in.JSONBytes()

			require.Nil(t, err)
			assert.Equal(t, []byte(tc.want), got)

		})
	}
}
