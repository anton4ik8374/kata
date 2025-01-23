package main

import (
	"reflect"
	"testing"
)

func TestNewProxy(t *testing.T) {
	type args struct {
		targetURL string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{name: "Proxy 1", args: args{targetURL: "http://127.0.0.1:8080"}, want: 200},
		{name: "Proxy 2", args: args{targetURL: "http://127.0.0.1:1313"}, want: 200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProxy(tt.args.targetURL); reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProxy() = %v, want %v", got, tt.want)
			}
		})
	}
}
