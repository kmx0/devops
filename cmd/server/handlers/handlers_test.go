package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kmx0/devops/cmd/server/storage"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleUpdate(t *testing.T) {
	SetRepository(storage.NewInMemory())
	type wantStruct struct {
		statusCode int
		// counter     types.Counter
	}
	// var store repositories.Repository

	router := SetupRouter()
	tests := []struct {
		name string
		req  string
		want wantStruct
	}{
		{
			name: "update",
			req:  "/update/counter/testCounter/100",
			want: wantStruct{
				statusCode: 200,
			},
		},
		{
			name: "without_id_counter",
			req:  "/update/counter/",
			want: wantStruct{
				statusCode: 404,
			},
		},
		{
			name: "invalid_value",
			req:  "/update/counter/testCounter/none",
			want: wantStruct{
				statusCode: 400,
			},
		},
		{
			name: "without_id_gauge",
			req:  "/update/gauge/",
			want: wantStruct{
				statusCode: 404,
			},
		},
		{
			name: "update_invalid_type",
			req:  "/update/gauge/testGauge/none",
			want: wantStruct{
				statusCode: 400,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			logrus.Info(tt.req)
			w := httptest.NewRecorder()
			// req, _ := http.NewRequest("GET", "/ping", nil)
			request, _ := http.NewRequest(http.MethodPost, tt.req, nil)

			router.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			err := res.Body.Close()
			require.NoError(t, err)
			// mapresult, err := ioutil.ReadAll(res.Body)
			// HandleCounter(tt.args.w, tt.args.r)
		})
	}
}

// func TestHandleUnknown(t *testing.T) {
// 	SetRepository(storage.NewInMemory())
// 	type wantStruct struct {
// 		contetnType string
// 		statusCode  int
// 		// counter     types.Counter
// 		// inmemWant repositories.Repository
// 	}
// 	// var store repositories.Repository

// 	tests := []struct {
// 		name     string
// 		req      string
// 		inmemReq repositories.Repository
// 		want     wantStruct
// 	}{
// 		{
// 			name:     "update_invalid_type",
// 			req:      "/update/unknown/testCounter/100",
// 			inmemReq: storage.NewInMemory(),
// 			want: wantStruct{
// 				statusCode:  501,
// 				contetnType: "",
// 				// inmemWant: &storage.InMemory{
// 				// 	MapCounter: make(map[string]types.Counter),
// 				// 	MapGauge:   make(map[string]types.Gauge),
// 				// },
// 			},
// 		},

// 		// {
// 		// 	name:     "test 3",
// 		// 	req:      "/update/counter/PollCount/1.1",
// 		// 	inmemReq: storage.NewInMemory(),
// 		// 	want: wantStruct{
// 		// 		statusCode:  500,
// 		// 		contetnType: "",
// 		// 		inmemWant: &storage.InMemory{
// 		// 			MapCounter: make(map[string]types.Counter),
// 		// 			MapGauge:   make(map[string]types.Gauge),
// 		// 		},
// 		// 	},
// 		// },
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			request := httptest.NewRequest(http.MethodGet, tt.req, nil)
// 			w := httptest.NewRecorder()

// 			h := http.HandlerFunc(HandleUnknown)
// 			h.ServeHTTP(w, request)
// 			res := w.Result()

// 			assert.Equal(t, tt.want.statusCode, res.StatusCode)
// 			assert.Equal(t, tt.want.contetnType, res.Header.Get("Content-Type"))
// 			// assert.Equal(t, tt.want.inmemWant, tt.inmemReq)
// 			err := res.Body.Close()
// 			require.NoError(t, err)
// 			// mapresult, err := ioutil.ReadAll(res.Body)
// 			// HandleCounter(tt.args.w, tt.args.r)
// 		})
// 	}
// }

func TestHandleValue(t *testing.T) {
	SetRepository(storage.NewInMemory())
	type wantStruct struct {
		statusCode int
		// counter     types.Counter
	}
	// var store repositories.Repository

	router := SetupRouter()
	tests := []struct {
		name string
		req  string
		want wantStruct
	}{
		// {
		// 	name: "update",
		// 	req:  "/update/counter/testCounter/100",
		// 	want: wantStruct{
		// 		statusCode: 200,
		// 	},
		// },
		// {
		// 	name: "without_id_counter",
		// 	req:  "/update/counter/",
		// 	want: wantStruct{
		// 		statusCode: 404,
		// 	},
		// },
		// {
		// 	name: "invalid_value",
		// 	req:  "/update/counter/testCounter/none",
		// 	want: wantStruct{
		// 		statusCode: 400,
		// 	},
		// },
		{
			name: "without_id_gauge",
			req:  "/update/gauge/",
			want: wantStruct{
				statusCode: 404,
			},
		},
		{
			name: "update_sequence",
			req:  "/value/gauge/testSetGet134",
			want: wantStruct{
				statusCode: 400,
			},
		},
	}
	w := httptest.NewRecorder()
	// requestPost, _ := http.NewRequest(http.MethodGet, "/value/gauge/testSetGet134", nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			logrus.Info(tt.req)
			request, _ := http.NewRequest(http.MethodGet, tt.req, nil)

			router.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			err := res.Body.Close()
			require.NoError(t, err)
			// mapresult, err := ioutil.ReadAll(res.Body)
			// HandleCounter(tt.args.w, tt.args.r)
		})
	}
}
