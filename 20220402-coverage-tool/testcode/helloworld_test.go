package testcode

import (
	"math/rand"
	"testing"
)

func Test_Hello(t *testing.T) {
	type args struct {
		name string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"normal",
			args{
				name: "world",
			},
			"hello world:10",
		},
	}
	for _, tt := range tests {
		Intn = MockIntn
		t.Run(tt.name, func(t *testing.T) {
			if got := Hello(tt.args.name); got != tt.want {
				t.Errorf("hello() = %v, want %v", got, tt.want)
			}
		})
		Intn = rand.Intn
	}
}
