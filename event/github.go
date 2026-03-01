package event

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// ErrNoEventPath is returned when GITHUB_EVENT_PATH is not set.
var ErrNoEventPath = errors.New("GITHUB_EVENT_PATH environment variable is not set")

// ErrNoCommitMessage is returned when the event JSON has no head_commit.message.
var ErrNoCommitMessage = errors.New("head_commit.message is empty or absent in event JSON (is this a push event?)")

type githubEvent struct {
	HeadCommit *struct {
		Message string `json:"message"`
	} `json:"head_commit"`
}

// ReadCommitMessage reads the commit message from the GitHub event JSON file
// pointed to by the GITHUB_EVENT_PATH environment variable.
func ReadCommitMessage() (string, error) {
	path := os.Getenv("GITHUB_EVENT_PATH")
	if path == "" {
		return "", ErrNoEventPath
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading event file %q: %w", path, err)
	}

	var event githubEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return "", fmt.Errorf("parsing event JSON: %w", err)
	}

	if event.HeadCommit == nil || event.HeadCommit.Message == "" {
		return "", ErrNoCommitMessage
	}

	return event.HeadCommit.Message, nil
}
