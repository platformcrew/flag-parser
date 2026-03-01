package output

import (
	"errors"
	"fmt"
	"os"

	"github.com/platformcrew/flag-parser/parser"
)

// ErrNoGithubOutput is returned when GITHUB_OUTPUT is not set.
var ErrNoGithubOutput = errors.New("GITHUB_OUTPUT environment variable is not set")

// WriteOutputs appends each flag result to the file pointed to by $GITHUB_OUTPUT.
func WriteOutputs(results []parser.FlagResult) error {
	path := os.Getenv("GITHUB_OUTPUT")
	if path == "" {
		return ErrNoGithubOutput
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("opening GITHUB_OUTPUT file %q: %w", path, err)
	}
	defer f.Close()

	for _, r := range results {
		value := "false"
		if r.Found {
			value = "true"
		}
		if _, err := fmt.Fprintf(f, "%s=%s\n", r.Name, value); err != nil {
			return fmt.Errorf("writing output for flag %q: %w", r.Name, err)
		}
	}

	return nil
}
