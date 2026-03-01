package parser

import (
	"testing"
)

func TestParseFlagDefinitions(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []FlagDef
		wantErr bool
	}{
		{
			name:  "single valid flag",
			input: `example-flag: "[example-flag]"`,
			want:  []FlagDef{{Name: "example-flag", SearchString: "[example-flag]"}},
		},
		{
			name: "multiple valid flags",
			input: `example-flag: "[example-flag]"
skip-tests: "[skip-tests]"
deploy-prod: "DEPLOY_PROD"`,
			want: []FlagDef{
				{Name: "example-flag", SearchString: "[example-flag]"},
				{Name: "skip-tests", SearchString: "[skip-tests]"},
				{Name: "deploy-prod", SearchString: "DEPLOY_PROD"},
			},
		},
		{
			name: "blank lines and comments skipped",
			input: `# this is a comment
example-flag: "[example-flag]"

# another comment
skip-tests: "[skip-tests]"`,
			want: []FlagDef{
				{Name: "example-flag", SearchString: "[example-flag]"},
				{Name: "skip-tests", SearchString: "[skip-tests]"},
			},
		},
		{
			name:  "windows line endings CRLF",
			input: "example-flag: \"[example-flag]\"\r\nskip-tests: \"[skip-tests]\"",
			want: []FlagDef{
				{Name: "example-flag", SearchString: "[example-flag]"},
				{Name: "skip-tests", SearchString: "[skip-tests]"},
			},
		},
		{
			name:  "colon inside search string",
			input: `flag-with-colon: "[search:colon]"`,
			want:  []FlagDef{{Name: "flag-with-colon", SearchString: "[search:colon]"}},
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			input:   "   \n  \n",
			wantErr: true,
		},
		{
			name:    "missing colon separator",
			input:   `example-flag "[example-flag]"`,
			wantErr: true,
		},
		{
			name:    "unquoted value",
			input:   `example-flag: [example-flag]`,
			wantErr: true,
		},
		{
			name:    "only opening quote",
			input:   `example-flag: "[example-flag]`,
			wantErr: true,
		},
		{
			name:    "empty quoted value",
			input:   `example-flag: ""`,
			wantErr: true,
		},
		{
			name:    "duplicate flag name",
			input:   "example-flag: \"[example-flag]\"\nexample-flag: \"[other]\"",
			wantErr: true,
		},
		{
			name:    "invalid flag name starts with number",
			input:   `1flag: "[test]"`,
			wantErr: true,
		},
		{
			name:    "empty flag name",
			input:   `: "[test]"`,
			wantErr: true,
		},
		{
			name:  "flag name with underscores and digits",
			input: `flag_name123: "[test]"`,
			want:  []FlagDef{{Name: "flag_name123", SearchString: "[test]"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFlagDefinitions(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseFlagDefinitions() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if len(got) != len(tt.want) {
				t.Fatalf("ParseFlagDefinitions() got %d defs, want %d", len(got), len(tt.want))
			}
			for i, def := range got {
				if def.Name != tt.want[i].Name || def.SearchString != tt.want[i].SearchString {
					t.Errorf("def[%d] = %+v, want %+v", i, def, tt.want[i])
				}
			}
		})
	}
}

func TestEvaluateFlags(t *testing.T) {
	defs := []FlagDef{
		{Name: "example-flag", SearchString: "[example-flag]"},
		{Name: "skip-tests", SearchString: "[skip-tests]"},
	}

	tests := []struct {
		name  string
		text  string
		found []bool
	}{
		{
			name:  "both flags found",
			text:  "deploy [example-flag] and [skip-tests] today",
			found: []bool{true, true},
		},
		{
			name:  "one flag found",
			text:  "build with [example-flag]",
			found: []bool{true, false},
		},
		{
			name:  "no flags found",
			text:  "regular commit message",
			found: []bool{false, false},
		},
		{
			name:  "empty text",
			text:  "",
			found: []bool{false, false},
		},
		{
			name:  "case sensitive: wrong case not found",
			text:  "[USE-DEPOT-RUNNER]",
			found: []bool{false, false},
		},
		{
			name:  "substring match works",
			text:  "prefix[example-flag]suffix",
			found: []bool{true, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := EvaluateFlags(defs, tt.text)
			if len(results) != len(defs) {
				t.Fatalf("EvaluateFlags() returned %d results, want %d", len(results), len(defs))
			}
			for i, r := range results {
				if r.Found != tt.found[i] {
					t.Errorf("results[%d].Found = %v, want %v (flag %q, search %q, text %q)",
						i, r.Found, tt.found[i], r.Name, r.SearchString, tt.text)
				}
			}
		})
	}
}
