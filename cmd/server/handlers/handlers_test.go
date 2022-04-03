package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kmx0/devops/cmd/server/storage"
	"github.com/kmx0/devops/internal/repositories"
	"github.com/kmx0/devops/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleCounter(t *testing.T) {
	SetRepository(storage.NewInMemory())
	type wantStruct struct {
		contetnType string
		statusCode  int
		// counter     types.Counter
		inmemWant repositories.Repository
	}
	// var store repositories.Repository

	tests := []struct {
		name     string
		req      string
		inmemReq repositories.Repository
		want     wantStruct
	}{
		// {
		// 	name:     "test 1",
		// 	req:      "/update/counter/PollCount/4",
		// 	inmemReq: storage.NewInMemory(),
		// 	want: wantStruct{
		// 		statusCode:  200,
		// 		contetnType: "text/plain",
		// 		inmemWant: &storage.InMemory{
		// 			MapCounter: make(map[string]types.Counter),
		// 			MapGauge:   make(map[string]types.Gauge),
		// 		},
		// 	},
		// },
		{
			name:     "without_id",
			req:      "/update/counter/",
			inmemReq: storage.NewInMemory(),
			want: wantStruct{
				statusCode:  404,
				contetnType: "",
				inmemWant: &storage.InMemory{
					MapCounter: make(map[string]types.Counter),
					MapGauge:   make(map[string]types.Gauge),
				},
			},
		},
		{
			name:     "invalid_value",
			req:      "/update/counter/testCounter/none",
			inmemReq: storage.NewInMemory(),
			want: wantStruct{
				statusCode:  400,
				contetnType: "",
				inmemWant: &storage.InMemory{
					MapCounter: make(map[string]types.Counter),
					MapGauge:   make(map[string]types.Gauge),
				},
			},
		},
		// {
		// 	name:     "test 3",
		// 	req:      "/update/counter/PollCount/1.1",
		// 	inmemReq: storage.NewInMemory(),
		// 	want: wantStruct{
		// 		statusCode:  500,
		// 		contetnType: "",
		// 		inmemWant: &storage.InMemory{
		// 			MapCounter: make(map[string]types.Counter),
		// 			MapGauge:   make(map[string]types.Gauge),
		// 		},
		// 	},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodGet, tt.req, nil)
			w := httptest.NewRecorder()

			h := http.HandlerFunc(HandleCounter)
			h.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			assert.Equal(t, tt.want.contetnType, res.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.inmemWant, tt.inmemReq)
			err := res.Body.Close()
			require.NoError(t, err)
			// mapresult, err := ioutil.ReadAll(res.Body)
			// HandleCounter(tt.args.w, tt.args.r)
		})
	}
}

func TestHandleGauge(t *testing.T) {
	SetRepository(storage.NewInMemory())
	type wantStruct struct {
		contetnType string
		statusCode  int
		inmemWant   repositories.Repository
	}
	tests := []struct {
		name     string
		req      string
		inmemReq repositories.Repository
		want     wantStruct
	}{
		{
			name:     "gauge test 1",
			req:      "/update/gauge/Alloc/24534",
			inmemReq: storage.NewInMemory(),
			want: wantStruct{
				statusCode:  200,
				contetnType: "text/plain",
				inmemWant: &storage.InMemory{
					MapCounter: make(map[string]types.Counter),
					MapGauge:   make(map[string]types.Gauge),
				},
			},
		},
		{
			name:     "gauge test 2",
			req:      "/update/gauge/BuckHashSys/1213.2",
			inmemReq: storage.NewInMemory(),
			want: wantStruct{
				statusCode:  200,
				contetnType: "text/plain",
				inmemWant: &storage.InMemory{
					MapCounter: make(map[string]types.Counter),
					MapGauge:   make(map[string]types.Gauge),
				},
			},
		},
		{
			name:     "gauge test 3",
			req:      "/update/gauge/RandomValue/",
			inmemReq: storage.NewInMemory(),
			want: wantStruct{
				statusCode:  400,
				contetnType: "",
				inmemWant: &storage.InMemory{
					MapCounter: make(map[string]types.Counter),
					MapGauge:   make(map[string]types.Gauge),
				},
			},
		},
		{
			name:     "without_id",
			req:      "/update/gauge/",
			inmemReq: storage.NewInMemory(),
			want: wantStruct{
				statusCode:  404,
				contetnType: "",
				inmemWant: &storage.InMemory{
					MapCounter: make(map[string]types.Counter),
					MapGauge:   make(map[string]types.Gauge),
				},
			},
		},
		{
			name:     "invalid_value",
			req:      "/update/gauge/testGauge/none",
			inmemReq: storage.NewInMemory(),
			want: wantStruct{
				statusCode:  400,
				contetnType: "",
				inmemWant: &storage.InMemory{
					MapCounter: make(map[string]types.Counter),
					MapGauge:   make(map[string]types.Gauge),
				},
			},
		}, //update_invalid_type
		{
			name:     "update_invalid_type",
			req:      "/update/gauge/testGauge/none",
			inmemReq: storage.NewInMemory(),
			want: wantStruct{
				statusCode:  400,
				contetnType: "",
				inmemWant: &storage.InMemory{
					MapCounter: make(map[string]types.Counter),
					MapGauge:   make(map[string]types.Gauge),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodGet, tt.req, nil)
			w := httptest.NewRecorder()

			h := http.HandlerFunc(HandleGauge)
			h.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			assert.Equal(t, tt.want.contetnType, res.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.inmemWant, tt.inmemReq)
			// mapresult, err := ioutil.ReadAll(res.Body)
			// HandleCounter(tt.args.w, tt.args.r)
			err := res.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestHandleUnknown(t *testing.T) {
	SetRepository(storage.NewInMemory())
	type wantStruct struct {
		contetnType string
		statusCode  int
		// counter     types.Counter
		inmemWant repositories.Repository
	}
	// var store repositories.Repository

	tests := []struct {
		name     string
		req      string
		inmemReq repositories.Repository
		want     wantStruct
	}{
		{
			name:     "update_invalid_type",
			req:      "/update/unknown/testCounter/100",
			inmemReq: storage.NewInMemory(),
			want: wantStruct{
				statusCode:  501,
				contetnType: "",
				inmemWant: &storage.InMemory{
					MapCounter: make(map[string]types.Counter),
					MapGauge:   make(map[string]types.Gauge),
				},
			},
		},

		// {
		// 	name:     "test 3",
		// 	req:      "/update/counter/PollCount/1.1",
		// 	inmemReq: storage.NewInMemory(),
		// 	want: wantStruct{
		// 		statusCode:  500,
		// 		contetnType: "",
		// 		inmemWant: &storage.InMemory{
		// 			MapCounter: make(map[string]types.Counter),
		// 			MapGauge:   make(map[string]types.Gauge),
		// 		},
		// 	},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodGet, tt.req, nil)
			w := httptest.NewRecorder()

			h := http.HandlerFunc(HandleUnknown)
			h.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			assert.Equal(t, tt.want.contetnType, res.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.inmemWant, tt.inmemReq)
			err := res.Body.Close()
			require.NoError(t, err)
			// mapresult, err := ioutil.ReadAll(res.Body)
			// HandleCounter(tt.args.w, tt.args.r)
		})
	}
}