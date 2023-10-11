package main

import "testing"

func Test_handler(t *testing.T) {
	type args struct {
		n int
		x int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"0", args{0, 0}, 0},
		{"1", args{1, 0}, 0},
		{"2", args{2, 1}, 3},
		{"3", args{4, 0}, 4},
		{"4", args{7, 1}, 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handler(tt.args.n, tt.args.x); got != tt.want {
				t.Errorf("handler() = %v, want %v", got, tt.want)
			}
		})
	}
}
