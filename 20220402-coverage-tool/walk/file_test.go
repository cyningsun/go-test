package walk

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestWild(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// {
		// 	"path",
		// 	args{
		// 		path: "../testcode",
		// 	},
		// 	[]string{
		// 		"../testcode",
		// 	},
		// 	false,
		// },
		// {
		// 	"path",
		// 	args{
		// 		path: "../testcode/",
		// 	},
		// 	[]string{
		// 		"../testcode/",
		// 	},
		// 	false,
		// },
		{
			"wild",
			args{
				path: "../testcode/*",
			},
			[]string{
				"../testcode/coverage.out",
				"../testcode/examples",
				"../testcode/go.mod",
				"../testcode/helloworld.go",
				"../testcode/helloworld_test.go",
				"../testcode/mock_intn.go",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filepath.Glob(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Children() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Children() = %v, want %v", got, tt.want)
			}
		})
	}
}
