package main

import "testing"

func Test_sendMetrics(t *testing.T) {
	type args struct {
		rm RunMetrics
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
