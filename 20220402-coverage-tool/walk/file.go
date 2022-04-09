package walk

import (
	"os"
	"path/filepath"
)

func Ancestor(path string) ([]string, error) {
	ancestor := make([]string, 0)
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	for dir != "/" {
		dir = filepath.Dir(dir)
		ancestor = append(ancestor, dir)
	}
	return ancestor, err
}
