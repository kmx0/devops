package crypto

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kmx0/devops/internal/types"
	"gotest.tools/assert"
)

func TestCheckHash(t *testing.T) {
	var helperi int64 = 1
	var helperf float64 = 1
	tests := []struct {
		name   string
		key    string
		metric types.Metrics
		want   error
	}{
		{
			name: "Counter nil pointer",
			metric: types.Metrics{
				ID:    "id1",
				MType: "counter",
			},
			key:  "hashkey",
			want: errors.New("received nil pointer on Delta"),
		},
		{
			name: "Gauge nil pointer",
			metric: types.Metrics{
				ID:    "id1",
				MType: "gauge",
			},
			key:  "hashkey",
			want: errors.New("received nil pointer on Value"),
		},
		{
			name: "Gauge empty key",
			metric: types.Metrics{
				ID:    "id1",
				MType: "gauge",
			},
			key: "",
		},
		{
			name: "Counter empty key",
			metric: types.Metrics{
				ID:    "id1",
				MType: "counter",
			},
			key: "",
		},
		{
			name: "Counter incorrect hash",
			metric: types.Metrics{
				ID:    "id1",
				MType: "counter",
				Delta: &helperi,
				Hash:  "",
			},
			key:  "hashkey",
			want: errors.New("hash sum not matched"),
		},
		{
			name: "Gauge incorrect hash",
			metric: types.Metrics{
				ID:    "id1",
				MType: "gauge",
				Value: &helperf,
				Hash:  "",
			},
			key:  "hashkey",
			want: errors.New("hash sum not matched"),
		},
		{
			name: "Gauge correct hash",
			metric: types.Metrics{
				ID:    "id1",
				MType: "gauge",
				Value: &helperf,
				Hash:  Hash(fmt.Sprintf("%s:gauge:%f", "id1", helperf), "hashkey"),
			},
			key: "hashkey",
		},
		{
			name: "Counter correct hash",
			metric: types.Metrics{
				ID:    "id1",
				MType: "counter",
				Delta: &helperi,
				Hash:  Hash(fmt.Sprintf("%s:counter:%d", "id1", helperi), "hashkey"),
			},
			key: "hashkey",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckHash(tt.metric, tt.key)
			if tt.want != nil && err != nil {
				assert.Equal(t, tt.want.Error(), err.Error())
			}
			if tt.want == nil || err == nil {
				assert.Equal(t, tt.want, err)
			}

			// assert
			// require.Error(t, err)

		})
	}
}

func ExampleHash() {

	src := "srcExample"
	key := "keyhash"
	hash := Hash(src, key)

	fmt.Println(hash)

	// Output:
	// 745a29df765bce14ea0c849cc49db8794ebb7b5bb7fe058465176a4dee698cb4
}
