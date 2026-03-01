package main

import (
	"fmt"
	"os"

	"github.com/platformcrew/flag-parser/event"
	"github.com/platformcrew/flag-parser/output"
	"github.com/platformcrew/flag-parser/parser"
)

func main() {
	if err := run(os.Getenv); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(getenv func(string) string) error {
	rawFlags := getenv("INPUT_FLAGS")
	if rawFlags == "" {
		return fmt.Errorf("flags input is required but was not provided")
	}

	defs, err := parser.ParseFlagDefinitions(rawFlags)
	if err != nil {
		return fmt.Errorf("parsing flags input: %w", err)
	}

	text := getenv("INPUT_TEXT")
	if text == "" {
		text, err = event.ReadCommitMessage()
		if err != nil {
			return fmt.Errorf("text input is empty and could not read commit message: %w", err)
		}
	}

	results := parser.EvaluateFlags(defs, text)

	if err := output.WriteOutputs(results); err != nil {
		return fmt.Errorf("writing outputs: %w", err)
	}

	for _, r := range results {
		status := "false"
		if r.Found {
			status = "true"
		}
		fmt.Printf("flag %s=%s (searched for %q)\n", r.Name, status, r.SearchString)
	}

	return nil
}
