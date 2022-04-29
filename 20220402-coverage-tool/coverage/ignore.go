package coverage

import (
	"github.com/cyningsun/go-test/20220402-coverage-tool/config"
	"github.com/cyningsun/go-test/20220402-coverage-tool/wildmatch"
)

func isIgnore(path string) MatchResult {
	for _, r := range config.IgnorePattern() {
		if wildmatch.WildMatch(r, path, wildmatch.WM_PATHNAME) == 0 {
			return MatchResult{path, r, true}
		}
	}

	return MatchResult{"", "", false}
}

type MatchResult struct {
	Path    string
	Pattern string
	OK      bool
}

func IgnoreFiles(pathes []string) ([]MatchResult, error) {
	matches := make([]MatchResult, 0, len(pathes))
	for _, each := range pathes {
		for _, r := range config.IgnorePattern() {
			if wildmatch.WildMatch(r, each, wildmatch.WM_PATHNAME) == 0 {
				matches = append(matches, MatchResult{each, r, true})
				break
			}
		}
	}

	return matches, nil
}
