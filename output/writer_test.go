package output

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/platformcrew/flag-parser/parser"
)

func setEnv(t *testing.T, key, value string) {
	t.Helper()
	old, had := os.LookupEnv(key)
	os.Setenv(key, value)
	t.Cleanup(func() {
		if had {
			os.Setenv(key, old)
		} else {
			os.Unsetenv(key)
		}
	})
}

func unsetEnv(t *testing.T, key string) {
	t.Helper()
	old, had := os.LookupEnv(key)
	os.Unsetenv(key)
	t.Cleanup(func() {
		if had {
			os.Setenv(key, old)
		}
	})
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}
	return string(data)
}

func TestWriteOutputs(t *testing.T) {
	t.Run("GITHUB_OUTPUT not set", func(t *testing.T) {
		unsetEnv(t, "GITHUB_OUTPUT")
		err := WriteOutputs([]parser.FlagResult{{Name: "flag", Found: true}})
		if !errors.Is(err, ErrNoGithubOutput) {
			t.Fatalf("got error %v, want ErrNoGithubOutput", err)
		}
	})

	t.Run("single flag found=true", func(t *testing.T) {
		f, _ := os.CreateTemp(t.TempDir(), "gh-output-*")
		f.Close()
		setEnv(t, "GITHUB_OUTPUT", f.Name())

		err := WriteOutputs([]parser.FlagResult{{Name: "example-flag", Found: true, SearchString: "[example-flag]"}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := readFile(t, f.Name())
		if got != "example-flag=true\n" {
			t.Errorf("got %q, want %q", got, "example-flag=true\n")
		}
	})

	t.Run("single flag found=false", func(t *testing.T) {
		f, _ := os.CreateTemp(t.TempDir(), "gh-output-*")
		f.Close()
		setEnv(t, "GITHUB_OUTPUT", f.Name())

		err := WriteOutputs([]parser.FlagResult{{Name: "skip-tests", Found: false, SearchString: "[skip-tests]"}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := readFile(t, f.Name())
		if got != "skip-tests=false\n" {
			t.Errorf("got %q, want %q", got, "skip-tests=false\n")
		}
	})

	t.Run("multiple flags written in order", func(t *testing.T) {
		f, _ := os.CreateTemp(t.TempDir(), "gh-output-*")
		f.Close()
		setEnv(t, "GITHUB_OUTPUT", f.Name())

		results := []parser.FlagResult{
			{Name: "example-flag", Found: true},
			{Name: "skip-tests", Found: false},
			{Name: "deploy-prod", Found: true},
		}
		if err := WriteOutputs(results); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := readFile(t, f.Name())
		want := "example-flag=true\nskip-tests=false\ndeploy-prod=true\n"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("appends to existing content", func(t *testing.T) {
		f, _ := os.CreateTemp(t.TempDir(), "gh-output-*")
		f.WriteString("existing=value\n")
		f.Close()
		setEnv(t, "GITHUB_OUTPUT", f.Name())

		if err := WriteOutputs([]parser.FlagResult{{Name: "new-flag", Found: true}}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got := readFile(t, f.Name())
		if !strings.HasPrefix(got, "existing=value\n") {
			t.Errorf("existing content was overwritten: %q", got)
		}
		if !strings.Contains(got, "new-flag=true\n") {
			t.Errorf("new flag not appended: %q", got)
		}
	})
}
