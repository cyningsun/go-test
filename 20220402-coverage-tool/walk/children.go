package walk

import (
	"os"
	"path/filepath"
)

func DescendantFile(path string) ([]string, error) {
	files := make([]string, 0)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	return files, err
}

func AncestorDir(path string) ([]string, error) {
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
