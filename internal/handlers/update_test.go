package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/viking311/monitoring/internal/entity"
)

func TestUpdateHandler_ServeHTTP(t *testing.T) {
	ch := make(chan entity.MetricEntityInterface, 100)
	handlerClass := NewUpdateHandler(ch)
	handler := http.Handler(handlerClass)
	type args struct {
		r *http.Request
	}

	tests := []struct {
		name string
		uh   UpdateHandler
		args args
		want int
	}{
		{
			name: "TestUpdateHandler_ServeHTTP_test1",
			args: args{

				r: httptest.NewRequest(http.MethodGet, "http://test.com/update/counter/1", nil),
			},
			want: http.StatusBadRequest,
		},
		{
			name: "TestUpdateHandler_ServeHTTP_test2",
			args: args{

				r: httptest.NewRequest(http.MethodPost, "http://test.com/update/counter/1", nil),
			},
			want: http.StatusNotFound,
		},
		{
			name: "TestUpdateHandler_ServeHTTP_test3",
			args: args{

				r: httptest.NewRequest(http.MethodPost, "http://test.com/update/counter/m1/1.8", nil),
			},
			want: http.StatusBadRequest,
		},
		{
			name: "TestUpdateHandler_ServeHTTP_test4",
			args: args{

				r: httptest.NewRequest(http.MethodPost, "http://test.com/update/gauge/m1/none", nil),
			},
			want: http.StatusBadRequest,
		},
		{
			name: "TestUpdateHandler_ServeHTTP_test5",
			args: args{

				r: httptest.NewRequest(http.MethodPost, "http://test.com/update/gauge/m1/1.8", nil),
			},
			want: http.StatusOK,
		},
		{
			name: "TestUpdateHandler_ServeHTTP_test6",
			args: args{

				r: httptest.NewRequest(http.MethodPost, "http://test.com/update/counter/m2/1", nil),
			},
			want: http.StatusOK,
		},
		{
			name: "TestUpdateHandler_ServeHTTP_test6",
			args: args{

				r: httptest.NewRequest(http.MethodPost, "http://test.com/update/gauge-new/m2/1", nil),
			},
			want: http.StatusNotImplemented,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, tt.args.r)
			assert.Equal(t, tt.want, rr.Code)
			if rr.Code == http.StatusOK {
				<-ch
			}
		})
	}
}
