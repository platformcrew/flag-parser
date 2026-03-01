package parser

import (
	"fmt"
	"regexp"
	"strings"
)

var validKeyRe = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9\-_]*$`)

// FlagDef holds a flag name and the literal string to search for.
type FlagDef struct {
	Name         string
	SearchString string
}

// FlagResult holds the evaluation result for a single flag.
type FlagResult struct {
	Name         string
	Found        bool
	SearchString string
}

// ParseFlagDefinitions parses a multiline string of flag definitions.
//
// Each line must be in the format:
//
//	flag-name: "search-string"
//
// Lines starting with '#' and blank lines are ignored.
func ParseFlagDefinitions(raw string) ([]FlagDef, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, fmt.Errorf("flags input is empty")
	}

	seen := make(map[string]bool)
	var defs []FlagDef

	for i, line := range strings.Split(raw, "\n") {
		lineNum := i + 1
		line = strings.TrimRight(line, "\r")
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		colonIdx := strings.Index(line, ":")
		if colonIdx < 0 {
			return nil, fmt.Errorf("line %d: missing ':' separator: %q", lineNum, line)
		}

		key := strings.TrimSpace(line[:colonIdx])
		value := strings.TrimSpace(line[colonIdx+1:])

		if key == "" {
			return nil, fmt.Errorf("line %d: flag name is empty", lineNum)
		}
		if !validKeyRe.MatchString(key) {
			return nil, fmt.Errorf("line %d: invalid flag name %q (must match ^[a-zA-Z][a-zA-Z0-9-_]*$)", lineNum, key)
		}
		if seen[key] {
			return nil, fmt.Errorf("line %d: duplicate flag name %q", lineNum, key)
		}

		if len(value) < 2 || value[0] != '"' || value[len(value)-1] != '"' {
			return nil, fmt.Errorf("line %d: search string for %q must be enclosed in double quotes, got: %s", lineNum, key, value)
		}
		searchString := value[1 : len(value)-1]
		if searchString == "" {
			return nil, fmt.Errorf("line %d: search string for %q must not be empty", lineNum, key)
		}

		seen[key] = true
		defs = append(defs, FlagDef{Name: key, SearchString: searchString})
	}

	if len(defs) == 0 {
		return nil, fmt.Errorf("flags input contains no valid flag definitions")
	}

	return defs, nil
}

// EvaluateFlags checks each flag definition against text using literal substring matching.
func EvaluateFlags(defs []FlagDef, text string) []FlagResult {
	results := make([]FlagResult, len(defs))
	for i, def := range defs {
		results[i] = FlagResult{
			Name:         def.Name,
			Found:        strings.Contains(text, def.SearchString),
			SearchString: def.SearchString,
		}
	}
	return results
}
