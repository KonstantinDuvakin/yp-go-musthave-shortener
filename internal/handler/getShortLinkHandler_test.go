package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KonstantinDuvakin/yp-go-musthave-shortener/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetShortLinkHandler(t *testing.T) {
	type args struct {
		storage *repository.LinksStorage
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
		method string
		args   args
		id     string
		want   want
	}{
		{
			name:   "Success redirect",
			method: http.MethodGet,
			args: args{
				storage: &repository.LinksStorage{
					Links: map[string]string{
						"1": "https://google.com",
					},
				},
			},
			id: "1",
			want: want{
				code: http.StatusTemporaryRedirect,
				body: "",
				headers: headers{
					{"Location", "https://google.com"},
				},
			},
		},
		{
			name:   "Not Found",
			method: http.MethodGet,
			args: args{
				storage: &repository.LinksStorage{
					Links: map[string]string{
						"1": "https://google.com",
					},
				},
			},
			id: "2",
			want: want{
				code:    http.StatusNotFound,
				body:    "Not Found",
				headers: headers{},
			},
		},
		{
			name:   "Method Not Allowed",
			method: http.MethodPost,
			args: args{
				storage: &repository.LinksStorage{
					Links: map[string]string{
						"1": "https://google.com",
					},
				},
			},
			id: "2",
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
			mux.Handle("GET /{id}", GetShortLinkHandler(tt.args.storage))

			request := httptest.NewRequest(tt.method, "/"+tt.id, nil)

			w := httptest.NewRecorder()

			mux.ServeHTTP(w, request)

			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)

			for _, h := range tt.want.headers {
				assert.Equal(t, h.value, res.Header.Get(h.key))
			}

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)

			assert.Equal(t, tt.want.body, string(resBody))
		})
	}
}
