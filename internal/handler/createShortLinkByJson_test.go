package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/KonstantinDuvakin/yp-go-musthave-shortener/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateShortLinkByJson(t *testing.T) {
	type args struct {
		s       *repository.LinksStorage
		baseUrl string
	}
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
		args   args
		method string
		body   io.Reader
		want   want
	}{
		{
			name: "CreateShortLinkByJson success",
			args: args{
				s: &repository.LinksStorage{
					Links:         make(map[string]string),
					RevertedLinks: make(map[string]string),
				},
				baseUrl: "http://localhost:8080",
			},
			method: http.MethodPost,
			body:   strings.NewReader(`{"url":"http://example.com"}`),
			want: want{
				code: http.StatusCreated,
				body: "http://localhost:8080/",
				headers: headers{
					{"Content-Type", "application/json"},
				},
			},
		},
		{
			name:   "CreateShortLinkHandler empty body",
			method: http.MethodPost,
			body:   nil,
			args: args{
				s:       &repository.LinksStorage{},
				baseUrl: "http://localhost:8080",
			},
			want: want{
				code:    http.StatusBadRequest,
				body:    "url is required",
				headers: headers{},
			},
		},
		{
			name:   "invalid JSON body",
			method: http.MethodPost,
			body:   strings.NewReader(`{"url":`), // оборванный JSON
			args: args{
				s: &repository.LinksStorage{
					Links:         make(map[string]string),
					RevertedLinks: make(map[string]string),
				},
				baseUrl: "http://localhost:8080",
			},
			want: want{
				code:    http.StatusBadRequest,
				body:    "invalid request body",
				headers: headers{},
			},
		},
		{
			name:   "valid JSON without url field",
			method: http.MethodPost,
			body:   strings.NewReader(`{}`), // декодится без ошибки, url == ""
			args: args{
				s: &repository.LinksStorage{
					Links:         make(map[string]string),
					RevertedLinks: make(map[string]string),
				},
				baseUrl: "http://localhost:8080",
			},
			want: want{
				code:    http.StatusBadRequest,
				body:    "url is required",
				headers: headers{},
			},
		},
		{
			name:   "empty url value",
			method: http.MethodPost,
			body:   strings.NewReader(`{"url":""}`),
			args: args{
				s: &repository.LinksStorage{
					Links:         make(map[string]string),
					RevertedLinks: make(map[string]string),
				},
				baseUrl: "http://localhost:8080",
			},
			want: want{
				code:    http.StatusBadRequest,
				body:    "url is required",
				headers: headers{},
			},
		},
		{
			name:   "unknown extra fields are ignored",
			method: http.MethodPost,
			body:   strings.NewReader(`{"url":"http://example.com","extra":123}`),
			args: args{
				s: &repository.LinksStorage{
					Links:         make(map[string]string),
					RevertedLinks: make(map[string]string),
				},
				baseUrl: "http://localhost:8080",
			},
			want: want{
				code: http.StatusCreated,
				body: "http://localhost:8080/",
				headers: headers{
					{"Content-Type", "application/json"},
				},
			},
		},
		{
			name:   "CreateShortLinkHandler not allowed method",
			method: http.MethodGet,
			body:   nil,
			args: args{
				s:       &repository.LinksStorage{},
				baseUrl: "http://localhost:8080",
			},
			want: want{
				code:    http.StatusMethodNotAllowed,
				body:    http.StatusText(http.StatusMethodNotAllowed),
				headers: headers{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.MethodNotAllowed(func(rw http.ResponseWriter, r *http.Request) {
				http.Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
			})

			r.Post("/", CreateShortLinkByJson(tt.args.s, tt.args.baseUrl))

			request := httptest.NewRequest(tt.method, "http://localhost:8080/", tt.body)

			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)

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
