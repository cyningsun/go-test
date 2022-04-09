package coverage

import (
	"github.com/cyningsun/go-test/20220402-coverage-tool/config"
	"github.com/cyningsun/go-test/20220402-coverage-tool/wildmatch"
)

func isIgnore(path string) bool {
	for _, r := range config.IgnorePattern() {
		if wildmatch.WildMatch(r, path, wildmatch.WM_PATHNAME) == 0 {
			return true
		}
	}

	return false
}

func IgnoreFiles(pathes []string) ([]string, error) {
	files := make([]string, 0, len(pathes))
	for _, each := range pathes {
		for _, r := range config.IgnorePattern() {
			if wildmatch.WildMatch(r, each, wildmatch.WM_PATHNAME) == 0 {
				files = append(files, each)
				break
			}
		}
	}

	return files, nil
}
