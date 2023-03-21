package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlers(t *testing.T) {
	todo := Todo{
		ID:    1,
		Title: "write a basic test",
	}

	tests := map[string]struct {
		in      http.HandlerFunc
		path    string
		method  string
		want    string
		body    string
		status  int
		muxvars map[string]string
		form    bool
	}{
		"healthz": {
			in:     healthHandler,
			method: http.MethodGet,
			status: http.StatusOK,
			path:   "/healthz",
			want:   `ok`,
		},
		"api/v1/healthz": {
			in:     healthHandler,
			method: http.MethodGet,
			status: http.StatusOK,
			path:   "/healthz",
			want:   `ok`,
		},
		"api/v1/todo": {
			in:     listHandler(NewMockStorage(false, false, todo)),
			method: http.MethodGet,
			status: http.StatusOK,
			path:   "/api/v1/todo",
			want:   `[{"id":1,"title":"write a basic test","updated":"0001-01-01T00:00:00Z","completed":"0001-01-01T00:00:00Z","complete":false}]`,
		},
		"api/v1/todo_get": {
			in:     readHandler(NewMockStorage(false, false, todo)),
			method: http.MethodGet,
			status: http.StatusOK,
			path:   "/api/v1/todo/1",
			want:   `{"id":1,"title":"write a basic test","updated":"0001-01-01T00:00:00Z","completed":"0001-01-01T00:00:00Z","complete":false}`,
			muxvars: map[string]string{
				"id": "1",
			},
		},
		"api/v1/todo_get_string_id": {
			in:     readHandler(NewMockStorage(false, false, todo)),
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			path:   "/api/v1/todo/test",
			want:   `{"text":"invalid! id must be integer","details":"todo id: test"}`,
			muxvars: map[string]string{
				"id": "test",
			},
		},
		"api/v1/todo_get_notexists": {
			in:     readHandler(NewMockStorage(false, false)),
			method: http.MethodGet,
			status: http.StatusNotFound,
			path:   "/api/v1/todo/2",
			want:   `{"text":"todo not found","details":"todo id: 2"}{"id":0,"title":"","updated":"0001-01-01T00:00:00Z","completed":"0001-01-01T00:00:00Z","complete":false}`,
			muxvars: map[string]string{
				"id": "2",
			},
		},
		"api/v1/todo_delete": {
			in:     deleteHandler(NewMockStorage(false, false, todo)),
			method: http.MethodDelete,
			status: http.StatusNoContent,
			path:   "/api/v1/todo/1",
			want:   `{"text":"todo deleted","details":"todo id: 1"}`,
			muxvars: map[string]string{
				"id": "1",
			},
		},

		"api/v1/todo_delete_error": {
			in:     deleteHandler(NewMockStorage(false, false, todo)),
			method: http.MethodDelete,
			status: http.StatusInternalServerError,
			path:   "/api/v1/todo/test",
			want:   `{"text":"invalid! id must be integer","details":"todo id: test"}`,
			muxvars: map[string]string{
				"id": "test",
			},
		},
		"api/v1/todo_create": {
			in:     createHandler(NewMockStorage(false, false)),
			method: http.MethodPost,
			status: http.StatusCreated,
			path:   "/api/v1/todo",
			want:   `{"id":1,"title":"write a basic test","updated":"0001-01-01T00:00:00Z","completed":"0001-01-01T00:00:00Z","complete":false}`,
			muxvars: map[string]string{
				"title": "write a basic test",
			},
			form: true,
		},
		"api/v1/todo_update": {
			in:     updateHandler(NewMockStorage(false, false, todo)),
			method: http.MethodPut,
			status: http.StatusOK,
			path:   "/api/v1/todo",
			want:   `{"id":1,"title":"write another test","updated":"0001-01-01T00:00:00Z","completed":"0001-01-01T00:00:00Z","complete":false}`,
			muxvars: map[string]string{
				"id":    "1",
				"title": "write another test",
			},
			form: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			body := []byte(tc.body)

			req := httptest.NewRequest(tc.method, tc.path, bytes.NewReader(body))
			w := httptest.NewRecorder()

			if tc.muxvars != nil {
				req = mux.SetURLVars(req, tc.muxvars)
			}

			if tc.form {
				urls := url.Values{}
				for i, v := range tc.muxvars {
					urls.Set(i, v)
				}
				req.Form = urls
			}

			tc.in(w, req)
			res := w.Result()
			defer res.Body.Close()
			got, err := ioutil.ReadAll(res.Body)

			require.Nil(t, err)
			assert.Equal(t, tc.want, string(got))
			assert.Equal(t, tc.status, w.Code)

		})
	}
}
