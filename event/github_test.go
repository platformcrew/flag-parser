package event

import (
	"errors"
	"os"
	"testing"
)

func writeEventFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "github-event-*.json")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

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

func TestReadCommitMessage(t *testing.T) {
	t.Run("GITHUB_EVENT_PATH not set", func(t *testing.T) {
		unsetEnv(t, "GITHUB_EVENT_PATH")
		_, err := ReadCommitMessage()
		if !errors.Is(err, ErrNoEventPath) {
			t.Fatalf("got error %v, want ErrNoEventPath", err)
		}
	})

	t.Run("valid event JSON with commit message", func(t *testing.T) {
		path := writeEventFile(t, `{"head_commit":{"message":"feat: add depot runner [use-depot-runner]"}}`)
		setEnv(t, "GITHUB_EVENT_PATH", path)
		msg, err := ReadCommitMessage()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := "feat: add depot runner [use-depot-runner]"
		if msg != want {
			t.Errorf("got %q, want %q", msg, want)
		}
	})

	t.Run("head_commit is null", func(t *testing.T) {
		path := writeEventFile(t, `{"head_commit":null}`)
		setEnv(t, "GITHUB_EVENT_PATH", path)
		_, err := ReadCommitMessage()
		if !errors.Is(err, ErrNoCommitMessage) {
			t.Fatalf("got error %v, want ErrNoCommitMessage", err)
		}
	})

	t.Run("head_commit absent", func(t *testing.T) {
		path := writeEventFile(t, `{"action":"opened"}`)
		setEnv(t, "GITHUB_EVENT_PATH", path)
		_, err := ReadCommitMessage()
		if !errors.Is(err, ErrNoCommitMessage) {
			t.Fatalf("got error %v, want ErrNoCommitMessage", err)
		}
	})

	t.Run("head_commit message empty", func(t *testing.T) {
		path := writeEventFile(t, `{"head_commit":{"message":""}}`)
		setEnv(t, "GITHUB_EVENT_PATH", path)
		_, err := ReadCommitMessage()
		if !errors.Is(err, ErrNoCommitMessage) {
			t.Fatalf("got error %v, want ErrNoCommitMessage", err)
		}
	})

	t.Run("malformed JSON", func(t *testing.T) {
		path := writeEventFile(t, `{not valid json`)
		setEnv(t, "GITHUB_EVENT_PATH", path)
		_, err := ReadCommitMessage()
		if err == nil {
			t.Fatal("expected error for malformed JSON, got nil")
		}
	})

	t.Run("file does not exist", func(t *testing.T) {
		setEnv(t, "GITHUB_EVENT_PATH", "/nonexistent/path/event.json")
		_, err := ReadCommitMessage()
		if err == nil {
			t.Fatal("expected error for missing file, got nil")
		}
	})
}
