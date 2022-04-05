package coverage

import (
	"regexp"

	"github.com/cyningsun/go-test/20220402-coverage-tool/config"
	"github.com/cyningsun/go-test/20220402-coverage-tool/wildmatch"
)

func Ignore(file string) bool {
	for _, r := range config.IgnorePattern() {
		if wildmatch.WildMatch(r, file, wildmatch.WM_PATHNAME) == 0 {
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
