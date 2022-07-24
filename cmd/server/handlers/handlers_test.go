package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kmx0/devops/cmd/server/storage"
	"github.com/kmx0/devops/internal/config"
	"github.com/kmx0/devops/internal/types"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleUpdate(t *testing.T) {
	s := storage.NewInMemory(config.Config{})
	SetRepository(s)
	type wantStruct struct {
		statusCode int
		// counter     types.Counter
	}
	// var store repositories.Repository

	router, _ := SetupRouter(config.Config{})
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

			// logrus.Info(tt.req)
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
func TestHandleUpdateJSON(t *testing.T) {
	s := storage.NewInMemory(config.Config{})
	SetRepository(s)
	type wantStruct struct {
		statusCode int
		// counter     types.Counter
	}

	router, _ := SetupRouter(config.Config{})
	tests := []struct {
		name string
		req  string
		body types.Metrics
		want wantStruct
	}{
		{
			name: "updateJSON_empty",
			req:  "/update/",
			body: types.Metrics{},
			want: wantStruct{
				statusCode: 501,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			bodyBytes, err := json.Marshal(tt.body)
			require.NoError(t, err)
			bodyReader := bytes.NewReader(bodyBytes)
			request, _ := http.NewRequest(http.MethodPost, tt.req, bodyReader)

			router.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			err = res.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestHandleValueJSON(t *testing.T) {
	s := storage.NewInMemory(config.Config{})
	SetRepository(s)
	type wantStruct struct {
		statusCode int
	}

	router, _ := SetupRouter(config.Config{})
	tests := []struct {
		name string
		req  string
		body types.Metrics
		want wantStruct
	}{
		{
			name: "valueJSON_POST_request_1",
			req:  "/value/",
			body: types.Metrics{
				ID:    "PollCount",
				MType: "counter",
			},
			want: wantStruct{
				statusCode: 404,
			},
		},
		{
			name: "valueJSON_POST_request_2",
			req:  "/value/",
			body: types.Metrics{
				ID:    "Alloc",
				MType: "gauge",
			},
			want: wantStruct{
				statusCode: 404,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			bodyBytes, err := json.Marshal(tt.body)
			require.NoError(t, err)
			bodyReader := bytes.NewReader(bodyBytes)
			request, _ := http.NewRequest(http.MethodPost, tt.req, bodyReader)

			router.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			err = res.Body.Close()
			require.NoError(t, err)
		})
	}
}

func ExampleHandlePing() {

	cfg.DBDSN = "postgres://postgres:postgres@localhost:5432/metrics"
	store = storage.NewInMemory(cfg)
	SetRepository(store)

	r := gin.Default()

	r.GET("/ping", HandlePing)

	// Listen and serve on 0.0.0.0:8080
	go r.Run(":8181")
	// Prepaire HTTP client
	time.Sleep(time.Second * 2)
	client := &http.Client{}

	endpoint := "http://127.0.0.1:8181/ping"
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		logrus.Error(err)
	}

	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	// печатаем код ответа
	fmt.Println(response.Status)
	// Output:
	// 200 OK
}

type testPostgres struct{}

func (t *testPostgres) Update(metric, name, value string) error {
	return nil
}
func (t *testPostgres) Set(ctx context.Context, key string, v interface{}) error {
	return nil
}

func (t *testPostgres) UpdateJSON(config.Config, types.Metrics) error {
	return nil
}
func (t *testPostgres) GetGauge(metric, name string) (g types.Gauge, e error) {
	if metric == "gauge" && name == "Alloc" {
		g = types.Gauge(1)
	} else {
		e = errors.New("not such metric")
	}
	return
}

func (t *testPostgres) GetCounterJSON(types.Metrics) (m types.Metrics, e error) {
	return
}
func (t *testPostgres) GetGaugeJSON(types.Metrics) (m types.Metrics, e error) {
	return
}

func (t *testPostgres) GetCounter(metric, name string) (c types.Counter, e error) {
	if metric == "counter" && name == "PollCount" {
		c = types.Counter(1)
	} else {
		e = errors.New("not such metric")
	}
	return
}
func (t *testPostgres) GetCurrentMetrics() (gm map[string]types.Gauge, cm map[string]types.Counter, e error) {
	gm = make(map[string]types.Gauge)
	cm = make(map[string]types.Counter)
	gm["Alloc"] = types.Gauge(1213)
	gm["Alloc2"] = types.Gauge(1214)
	cm["PollCount"] = types.Counter(12)
	cm["PollCount2"] = types.Counter(14)
	return
}
func (t *testPostgres) RestoreFromDisk(cfg config.Config) {}

func (t *testPostgres) SaveToDisk(cfg config.Config) {}

// func HandleAllValues(c *gin.Context) {
// 	mapGauge, mapCounter, _ := store.GetCurrentMetrics()

// 	c.Header("Content-Type", "text/html; charset=utf-8")
// 	c.String(http.StatusOK, "%+v\n%+v", mapGauge, mapCounter)
// }

func TestHandleAllValues(t *testing.T) {

	store = &testPostgres{}
	SetRepository(store)
	cfg = config.Config{}

	SetRepository(store)

	r := gin.New()
	r.Use(gin.Recovery(),
		Compress(),
		Decompress(),
		gin.Logger())

	r.GET("/", HandleAllValues)

	go r.Run(":8182")
	time.Sleep(time.Second * 2)
	client := &http.Client{}

	endpoint := "http://127.0.0.1:8182/"
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		logrus.Error(err)
	}

	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	payload, err := io.ReadAll(response.Body)

	gm := make(map[string]types.Gauge)
	cm := make(map[string]types.Counter)
	gm["Alloc"] = types.Gauge(1213)
	gm["Alloc2"] = types.Gauge(1214)
	cm["PollCount"] = types.Counter(12)
	cm["PollCount2"] = types.Counter(14)
	expect := fmt.Sprintf("%+v\n%+v", gm, cm)
	t.Run("1 test", func(t *testing.T) {

		assert.Equal(t, err, nil)
		assert.Equal(t, response.StatusCode, 200)
		assert.Equal(t, payload, []byte(expect))
	})

	//Checking

}

func TestHandleValue(t *testing.T) {
	// s := storage.NewInMemory(config.Config{})
	// SetRepository(s)

	store = &testPostgres{}
	SetRepository(store)
	cfg = config.Config{}

	r := gin.New()
	r.Use(gin.Recovery(),
		Compress(),
		Decompress(),
		gin.Logger())

	r.GET("/value/:typem/:metric", HandleValue)

	type wantStruct struct {
		statusCode int
	}
	// var store repositories.Repository

	tests := []struct {
		name string
		req  string
		want wantStruct
	}{
		{
			name: "Counter real value",
			req:  "/value/counter/PollCount",
			want: wantStruct{
				statusCode: 200,
			},
		},
		{
			name: "Gauge real value",
			req:  "/value/gauge/Alloc",
			want: wantStruct{
				statusCode: 200,
			},
		},
		{
			name: "Unknown Type",
			req:  "/value/unknowtype",
			want: wantStruct{
				statusCode: 404,
			},
		},
		{
			name: "Counter_invalid_metric_name",
			req:  "/value/counter/testCounter",
			want: wantStruct{
				statusCode: 404,
			},
		},
		{
			name: "Guage_invalid_metric_name",
			req:  "/value/gauge/testGauge",
			want: wantStruct{
				statusCode: 404,
			},
		},
	}
	// requestPost, _ := http.NewRequest(http.MethodGet, "/value/gauge/testSetGet134", nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			logrus.Info(tt.req)
			request, _ := http.NewRequest(http.MethodGet, tt.req, nil)

			r.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			err := res.Body.Close()
			require.NoError(t, err)
		})
	}
}
