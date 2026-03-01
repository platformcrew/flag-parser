# flag-parser

[![CI](https://github.com/platformcrew/flag-parser/actions/workflows/ci.yml/badge.svg)](https://github.com/platformcrew/flag-parser/actions/workflows/ci.yml)

A GitHub Action that parses boolean flags from text using user-defined key-value pattern matching.

## How it works

You define a list of flags, each mapping a flag name to a literal search string. The action checks whether each search string appears in the input text (case-sensitive substring match). If the `text` input is omitted, the action falls back to reading the commit message from the GitHub event context.

Each flag produces a step output with the value `"true"` or `"false"`.

## Usage

```yaml
- name: Parse flags
  id: flags
  uses: platformcrew/flag-parser@v1
  with:
    flags: |
      example-flag: "[example-flag]"
      skip-tests: "[skip-tests]"
      deploy-prod: "DEPLOY_PROD"
    # text: optional — defaults to the commit message on push events

- name: Check example-flag
  if: steps.flags.outputs.example-flag == 'true'
  run: echo "example-flag is enabled"

- name: Skip tests
  if: steps.flags.outputs.skip-tests == 'true'
  run: echo "Tests are skipped"
```

## Inputs

| Input   | Required | Default | Description                                                |
| ------- | -------- | ------- | ---------------------------------------------------------- |
| `flags` | Yes      | —       | Multiline flag definitions (see format below)              |
| `text`  | No       | `''`    | Text to search. Falls back to the commit message if empty. |

### `flags` format

Each non-blank, non-comment line must follow this format:

```
flag-name: "search-string"
```

- **flag-name**: Must start with a letter; may contain letters, digits, hyphens (`-`), and underscores (`_`).
- **search-string**: Must be enclosed in double quotes. The search is a literal, case-sensitive substring match — no regex, no globbing.
- Lines starting with `#` are treated as comments and ignored.

```yaml
flags: |
  # Enable example-flag runner when commit message contains [example-flag]
  example-flag: "[example-flag]"

  # Skip CI tests
  skip-tests: "[skip-tests]"
```

## Outputs

One output per flag, named after the flag. Value is always `"true"` or `"false"`.

```yaml
steps.flags.outputs.example-flag  # "true" or "false"
steps.flags.outputs.skip-tests    # "true" or "false"
```

## Commit message fallback

When `text` is not provided, the action reads `head_commit.message` from the GitHub event JSON (`GITHUB_EVENT_PATH`). This is only available on **push** events. For other event types (pull_request, workflow_dispatch, etc.), you must supply the `text` input explicitly.

Example using the PR body as the text source:

```yaml
- name: Parse flags from PR body
  id: flags
  uses: platformcrew/flag-parser@v1
  with:
    text: ${{ github.event.pull_request.body }}
    flags: |
      example-flag: "[example-flag]"
```

## Local smoke test

```bash
output=$(mktemp)
docker build -t flag-parser:test .
docker run --rm \
  -e INPUT_FLAGS='example-flag: "[example-flag]"' \
  -e INPUT_TEXT='build with [example-flag] enabled' \
  -e GITHUB_OUTPUT=/github/output \
  -v "$output:/github/output" \
  flag-parser:test
cat "$output"
# example-flag=true
```
