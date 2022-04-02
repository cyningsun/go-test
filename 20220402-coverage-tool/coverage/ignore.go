package coverage

import (
	"regexp"

	"github.com/cyningsun/go-test/20220402-coverage-tool/config"
)

func Ignore(file string) bool {
	for _, r := range config.IgnorePattern() {
		if r.MatchString(file) {
			return true
		}
	}

	return false
}

func CheckIgnore(pattern string, files []string) ([]string, error) {
	newFiles := make([]string, 0, len(files))
	r, _ := regexp.Compile(pattern)
	for _, each := range files {
		if !r.MatchString(each) {
			continue
		}

		newFiles = append(newFiles, each)
	}

	return newFiles, nil
}
