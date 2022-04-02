package coverage

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
)

type Line struct {
	FilePath string
	LineFrom int
	LineTo   int
	Count    int
	Covered  bool
}

type File struct {
	FilePath       string
	TotalLines     int
	CoveredLines   int
	UncoveredLines int
}

func (f *File) Coverage() float32 {
	if f.TotalLines == 0 {
		return 0
	}

	return 100 * float32(f.CoveredLines) / float32(f.TotalLines)
}

type Parser struct {
	Lines    []Line
	Files    []File
	TotalCov float32
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(cover string) error {
	raw, err := read(cover)
	if err != nil {
		return err
	}

	p.Lines = match(raw)
	p.Files = calFileCov(p.Lines)
	p.TotalCov = calTotalCov(p.Files)

	return nil
}

func read(cover string) ([]string, error) {
	file, err := os.Open(cover)
	if err != nil {
		return nil, err
	}

	bio := bufio.NewReader(file)
	lines := make([]string, 0)
	for {
		line, err := bio.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}

		if err == io.EOF {
			break
		}

		lines = append(lines, string(line))
	}
	return lines, nil
}

func match(input []string) []Line {
	re := regexp.MustCompile(`([^:]*):(\d+)\.\d*,(\d+)\.\d* (\d+) (\d+)`)

	output := make([]Line, 0, len(input))
	for _, each := range input {
		match := re.FindStringSubmatch(each)
		if len(match) != 6 {
			continue
		}

		filepath := match[1]
		from, _ := strconv.Atoi(match[2])
		to, _ := strconv.Atoi(match[3])
		count, _ := strconv.Atoi(match[4])
		covered, _ := strconv.ParseBool(match[5])

		if Ignore(filepath) {
			continue
		}

		output = append(output, Line{
			FilePath: filepath,
			LineFrom: from,
			LineTo:   to,
			Count:    count,
			Covered:  covered,
		})
	}

	return output
}

func calFileCov(lines []Line) []File {
	dict := make(map[string][]Line)
	for _, each := range lines {
		key := each.FilePath

		val, ok := dict[key]
		if !ok {
			val = make([]Line, 0)
		}

		val = append(val, each)
		dict[key] = val
	}

	files := make([]File, 0, len(dict))
	for filePath, file := range dict {
		sort.Slice(file, func(i, j int) bool {
			return file[i].LineTo < file[j].LineTo
		})

		total, covered, uncovered := 0, 0, 0
		for _, l := range file {
			total += l.Count
			if l.Covered {
				covered += l.Count
			} else {
				uncovered += l.Count
			}
		}

		files = append(files, File{
			FilePath:       filePath,
			TotalLines:     total,
			CoveredLines:   covered,
			UncoveredLines: uncovered,
		})
	}

	return files
}

func calTotalCov(files []File) float32 {
	total, covered := 0, 0
	for _, f := range files {
		covered += f.CoveredLines
		total += f.TotalLines
	}
	return 100 * float32(covered) / float32(total)
}

func (p *Parser) Println() {
	fmt.Println("totalCov:", p.TotalCov)
	for _, each := range p.Files {
		fmt.Println(each.FilePath, each.Coverage())
	}
}
