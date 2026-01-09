package main

import (
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Init kicks off the database connector
func (s *SQLStorage) MockInit() (*sqlmock.Sqlmock, error) {
	var err error
	var mock sqlmock.Sqlmock

	s.db, mock, err = sqlmock.New()
	if err != nil {
		return nil, err
	}

	return &mock, nil
}

func TestSQLCreate(t *testing.T) {
	tests := map[string]struct {
		in         Todo
		want       Todo
		expectExec string
		expectArgs []driver.Value
	}{
		"basic": {
			in: Todo{
				Title: "write a basic test",
			},
			want: Todo{
				ID:    1,
				Title: "write a basic test",
			},
			expectExec: `INSERT INTO todo`,
			expectArgs: []driver.Value{"write a basic test"},
		},
		"basicCompleted": {
			in: Todo{
				Title:    "write a basic test",
				Complete: true,
			},
			want: Todo{
				ID:       1,
				Title:    "write a basic test",
				Complete: true,
			},
			expectExec: `INSERT INTO todo`,
			expectArgs: []driver.Value{"write a basic test"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := SQLStorage{}
			m, err := s.MockInit()
			mock := *m
			require.Nil(t, err)

			mock.ExpectPrepare(tc.expectExec).
				ExpectExec().
				WithArgs(tc.expectArgs...).
				WillReturnResult(sqlmock.NewResult(1, 1))

			got, err := s.Create(tc.in)
			require.Nil(t, err)

			assert.Equal(t, tc.want, got)
			require.Nil(t, mock.ExpectationsWereMet())

		})
	}
}

func TestSQLRead(t *testing.T) {
	now := time.Now()
	tests := map[string]struct {
		in   *sqlmock.Rows
		id   string
		want Todo
	}{
		"basic": {
			in: sqlmock.NewRows([]string{"id", "title", "updated", "created"}).
				AddRow(1, "write a basic test", now, sql.NullTime{}),
			id: "1",
			want: Todo{
				ID:      1,
				Title:   "write a basic test",
				Updated: now,
				completedNull: sql.NullTime{
					Valid: false,
				},
			},
		},
		"basicCompleted": {
			in: sqlmock.NewRows([]string{"id", "title", "updated", "completed"}).
				AddRow(1, "write a basic test", now, now),
			id: "1",
			want: Todo{
				ID:        1,
				Title:     "write a basic test",
				Updated:   now,
				Completed: now,
				Complete:  true,
				completedNull: sql.NullTime{
					Time:  now,
					Valid: true,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := SQLStorage{}
			m, err := s.MockInit()
			mock := *m
			require.Nil(t, err)

			mock.ExpectQuery("SELECT (.+) FROM").WithArgs(tc.id).WillReturnRows(tc.in)

			got, err := s.Read(tc.id)
			require.Nil(t, err)

			assert.Equal(t, tc.want, got)

			require.Nil(t, mock.ExpectationsWereMet())

		})
	}
}

func TestSQLUpdate(t *testing.T) {
	now := time.Now()
	tests := map[string]struct {
		in         Todo
		id         string
		want       Todo
		expectExec string
		expectArgs []driver.Value
		selectin   *sqlmock.Rows
	}{
		"basic": {
			in: Todo{
				ID:    1,
				Title: "write a basic test",
			},
			want: Todo{
				ID:    1,
				Title: "write a basic test",
			},
			id: "1",
			selectin: sqlmock.NewRows([]string{"id", "title", "updated", "completed"}).
				AddRow(1, "write a basic test", now, sql.NullTime{}),
			expectExec: `UPDATE todo`,
			expectArgs: []driver.Value{"write a basic test", 1},
		},
		"basicCompleted": {
			in: Todo{
				ID:       1,
				Title:    "write a basic test",
				Complete: true,
				Updated:  now,
			},
			want: Todo{
				ID:        1,
				Title:     "write a basic test",
				Complete:  true,
				Updated:   now,
				Completed: now,
			},
			id: "1",
			selectin: sqlmock.NewRows([]string{"id", "title", "updated", "completed"}).
				AddRow(1, "write a basic test", now, sql.NullTime{}),
			expectExec: `UPDATE todo`,
			expectArgs: []driver.Value{"write a basic test", 1},
		},
		"basicInCompleted": {
			in: Todo{
				ID:        1,
				Title:     "write a basic test",
				Updated:   now,
				Completed: now,
			},
			want: Todo{
				ID:       1,
				Title:    "write a basic test",
				Complete: false,
				Updated:  now,
			},
			id: "1",
			selectin: sqlmock.NewRows([]string{"id", "title", "updated", "completed"}).
				AddRow(1, "write a basic test", now, now),
			expectExec: `UPDATE todo`,
			expectArgs: []driver.Value{"write a basic test", 1},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := SQLStorage{}
			m, err := s.MockInit()
			mock := *m
			require.Nil(t, err, "could not init mock")

			mock.ExpectQuery("SELECT (.+) FROM").WithArgs(tc.id).WillReturnRows(tc.selectin)

			mock.ExpectPrepare(tc.expectExec).
				ExpectExec().
				WithArgs(tc.expectArgs...).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err = s.Update(tc.in)
			require.Nil(t, err, "coult not update mock")

			require.Nil(t, mock.ExpectationsWereMet())

		})
	}
}

func TestSQLList(t *testing.T) {
	now := time.Now()
	tests := map[string]struct {
		in   *sqlmock.Rows
		want Todos
	}{
		"basic": {
			in: sqlmock.NewRows([]string{"id", "title", "updated", "created"}).
				AddRow(1, "write a basic test", now, sql.NullTime{}),
			want: Todos{
				Todo{
					ID:      1,
					Title:   "write a basic test",
					Updated: now,
					completedNull: sql.NullTime{
						Valid: false,
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := SQLStorage{}
			m, err := s.MockInit()
			mock := *m
			require.Nil(t, err)

			mock.ExpectQuery("SELECT (.+) FROM").WillReturnRows(tc.in)

			got, err := s.List()
			require.Nil(t, err)
			assert.Equal(t, tc.want, got)

			require.Nil(t, mock.ExpectationsWereMet())

		})
	}
}

func TestSQLDelete(t *testing.T) {
	tests := map[string]struct {
		id         string
		expectExec string
		expectArgs []driver.Value
	}{
		"basic": {
			id:         "1",
			expectExec: `DELETE FROM todo`,
			expectArgs: []driver.Value{"1"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := SQLStorage{}
			m, err := s.MockInit()
			mock := *m
			require.Nil(t, err)

			mock.ExpectPrepare(tc.expectExec).
				ExpectExec().
				WithArgs(tc.expectArgs...).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err = s.Delete(tc.id)
			require.Nil(t, err)

			require.Nil(t, mock.ExpectationsWereMet())

		})
	}
}
