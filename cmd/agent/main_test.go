package main

import (
	"testing"

	"github.com/kmx0/devops/internal/types"
)

func Test_sendMetrics(t *testing.T) {
	type args struct {
		rm types.RunMetrics
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendMetrics(tt.args.rm)
		})
	}
}
