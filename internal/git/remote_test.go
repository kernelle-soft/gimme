package git

import "testing"

func TestNormalizeRemoteURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// GitHub SSH format
		{
			name:     "github ssh with .git",
			input:    "git@github.com:user/repo.git",
			expected: "github.com/user/repo",
		},
		{
			name:     "github ssh without .git",
			input:    "git@github.com:user/repo",
			expected: "github.com/user/repo",
		},

		// GitHub HTTPS format
		{
			name:     "github https with .git",
			input:    "https://github.com/user/repo.git",
			expected: "github.com/user/repo",
		},
		{
			name:     "github https without .git",
			input:    "https://github.com/user/repo",
			expected: "github.com/user/repo",
		},

		// GitHub ssh:// format
		{
			name:     "github ssh:// with .git",
			input:    "ssh://git@github.com/user/repo.git",
			expected: "github.com/user/repo",
		},
		{
			name:     "github ssh:// without .git",
			input:    "ssh://git@github.com/user/repo",
			expected: "github.com/user/repo",
		},

		// Edge cases
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "whitespace only",
			input:    "   ",
			expected: "",
		},
		{
			name:     "input with whitespace",
			input:    "  git@github.com:user/repo.git  ",
			expected: "github.com/user/repo",
		},
		{
			name:     "unrecognized format",
			input:    "not-a-url",
			expected: "",
		},

		// HTTP (less common but valid)
		{
			name:     "http format",
			input:    "http://github.com/user/repo.git",
			expected: "github.com/user/repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeRemoteURL(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeRemoteURL(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
