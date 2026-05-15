package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Testcase struct {
	Index      int
	InputFile  string
	OutputFile string
}

type Problem struct {
	ID        string
	Dir       string
	Testcases []Testcase
}

func Scan(root string) ([]Problem, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("cannot read testcases directory %q: %w", root, err)
	}

	var problems []Problem
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		probDir := filepath.Join(root, entry.Name())
		tcs, err := scanTestcases(probDir)
		if err != nil {
			return nil, fmt.Errorf("problem %q: %w", entry.Name(), err)
		}
		if len(tcs) == 0 {
			continue
		}

		problems = append(problems, Problem{
			ID:        entry.Name(),
			Dir:       probDir,
			Testcases: tcs,
		})
	}

	sort.Slice(problems, func(i, j int) bool {
		ni, errI := strconv.Atoi(problems[i].ID)
		nj, errJ := strconv.Atoi(problems[j].ID)
		if errI == nil && errJ == nil {
			return ni < nj
		}
		return problems[i].ID < problems[j].ID
	})

	return problems, nil
}

func scanTestcases(dir string) ([]Testcase, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot read directory: %w", err)
	}

	inputFiles := make(map[string]string)
	outputFiles := make(map[string]string)

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".a") {
			key := strings.TrimSuffix(name, ".a")
			outputFiles[key] = filepath.Join(dir, name)
		} else {
			inputFiles[name] = filepath.Join(dir, name)
		}
	}

	var testcases []Testcase
	for key, inputPath := range inputFiles {
		outputPath, ok := outputFiles[key]
		if !ok {
			continue
		}

		idx, err := parseIndex(key)
		if err != nil {
			continue
		}

		testcases = append(testcases, Testcase{
			Index:      idx,
			InputFile:  inputPath,
			OutputFile: outputPath,
		})
	}

	sort.Slice(testcases, func(i, j int) bool {
		return testcases[i].Index < testcases[j].Index
	})

	return testcases, nil
}

func parseIndex(name string) (int, error) {
	n, err := strconv.Atoi(name)
	if err != nil {
		return 0, fmt.Errorf("cannot parse index from %q: %w", name, err)
	}
	return n, nil
}
