package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/KonstantinDuvakin/yp-go-musthave-shortener/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateShortLinkHandler(t *testing.T) {
	type headers []struct {
		key   string
		value string
	}
	type want struct {
		code    int
		body    string
		headers headers
	}
	tests := []struct {
		name   string
		method string
		body   io.Reader
		arg    *repository.LinksStorage
		want   want
	}{
		{
			name:   "CreateShortLinkHandler success creation",
			method: http.MethodPost,
			body:   strings.NewReader("http://link.com"),
			arg:    &repository.LinksStorage{},
			want: want{
				code: http.StatusCreated,
				body: "http://localhost:8080/",
				headers: headers{
					{"Content-Type", "text/plain"},
				},
			},
		},
		{
			name:   "CreateShortLinkHandler empty body",
			method: http.MethodPost,
			body:   nil,
			arg:    &repository.LinksStorage{},
			want: want{
				code:    http.StatusBadRequest,
				body:    "body is required",
				headers: headers{},
			},
		},
		{
			name:   "CreateShortLinkHandler not allowed method",
			method: http.MethodGet,
			body:   nil,
			arg:    &repository.LinksStorage{},
			want: want{
				code:    http.StatusMethodNotAllowed,
				body:    "Method Not Allowed\n",
				headers: headers{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.Handle("POST /", CreateShortLinkHandler(tt.arg))

			request := httptest.NewRequest(tt.method, "http://localhost:8080/", tt.body)

			w := httptest.NewRecorder()

			mux.ServeHTTP(w, request)

			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)
			for _, h := range tt.want.headers {
				assert.Equal(t, h.value, res.Header.Get(h.key))
			}

			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)

			require.NoError(t, err)

			assert.Contains(t, string(body), tt.want.body)
		})
	}
}
