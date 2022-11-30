package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/viking311/monitoring/internal/storage"
)

func TestUpdateHandler_ServeHTTP(t *testing.T) {
	// ch := make(chan entity.MetricEntityInterface, 100)
	s := storage.NewInMemoryStorage()
	handlerClass := NewUpdateHandler(s)
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", handlerClass.ServeHTTP)
	ts := httptest.NewServer(r)
	defer ts.Close()

	type args struct {
		url    string
		method string
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "TestUpdateHandler_ServeHTTP_test1",
			args: args{
				url:    ts.URL + "/update/counter/1",
				method: http.MethodGet,
			},
			want: http.StatusNotFound,
		},
		{
			name: "TestUpdateHandler_ServeHTTP_test2",
			args: args{
				url:    ts.URL + "/update/counter/1",
				method: http.MethodPost,
			},
			want: http.StatusNotFound,
		},
		{
			name: "TestUpdateHandler_ServeHTTP_test3",
			args: args{
				url:    ts.URL + "/update/counter/m1/1.8",
				method: http.MethodPost,
			},
			want: http.StatusBadRequest,
		},
		{
			name: "TestUpdateHandler_ServeHTTP_test4",
			args: args{
				url:    ts.URL + "/update/gauge/m1/none",
				method: http.MethodPost,
			},
			want: http.StatusBadRequest,
		},
		{
			name: "TestUpdateHandler_ServeHTTP_test5",
			args: args{
				url:    ts.URL + "/update/gauge/m1/1.8",
				method: http.MethodPost,
			},
			want: http.StatusOK,
		},
		{
			name: "TestUpdateHandler_ServeHTTP_test6",
			args: args{
				url:    ts.URL + "/update/counter/m2/1",
				method: http.MethodPost,
			},
			want: http.StatusOK,
		},
		{
			name: "TestUpdateHandler_ServeHTTP_test6",
			args: args{
				url:    ts.URL + "/update/gauge-new/m2/1",
				method: http.MethodPost,
			},
			want: http.StatusNotImplemented,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _ := sendTestRequest(t, tt.args.method, tt.args.url)
			assert.Equal(t, tt.want, code)
		})
	}
}
