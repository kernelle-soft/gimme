package git

import (
	"strings"
)

// NormalizeRemoteURL converts a git remote URL to a normalized identifier.
// Examples:
//   - git@github.com:user/repo.git     → github.com/user/repo
//   - git@github.com:user/repo         → github.com/user/repo
//   - https://github.com/user/repo.git → github.com/user/repo
//   - https://github.com/user/repo     → github.com/user/repo
//   - ssh://git@github.com/user/repo   → github.com/user/repo
//
// Returns empty string if the URL format is not recognized.
func NormalizeRemoteURL(url string) string {
	url = strings.TrimSpace(url)
	if url == "" {
		return ""
	}

	// Remove .git suffix if present
	url = strings.TrimSuffix(url, ".git")

	// Handle SSH format: git@github.com:user/repo
	if strings.HasPrefix(url, "git@") {
		// git@github.com:user/repo → github.com/user/repo
		url = strings.TrimPrefix(url, "git@")
		url = strings.Replace(url, ":", "/", 1)
		return url
	}

	// Handle ssh:// format: ssh://git@github.com/user/repo
	if strings.HasPrefix(url, "ssh://") {
		url = strings.TrimPrefix(url, "ssh://")
		url = strings.TrimPrefix(url, "git@")
		return url
	}

	// Handle https:// format: https://github.com/user/repo
	if strings.HasPrefix(url, "https://") {
		url = strings.TrimPrefix(url, "https://")
		return url
	}

	// Handle http:// format: http://github.com/user/repo
	if strings.HasPrefix(url, "http://") {
		url = strings.TrimPrefix(url, "http://")
		return url
	}

	// Unrecognized format
	return ""
}
