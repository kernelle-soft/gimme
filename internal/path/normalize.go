package path

import (
    "os"
    "path/filepath"
)

func Normalize(path string) (string, error) {
    // Expand environment variables first
    path = os.ExpandEnv(path)

    // Handle tilde
    if len(path) > 0 && path[0] == '~' {
        home, err := os.UserHomeDir()
        if err != nil {
            return "", err
        }
        path = home + path[1:]
    }

    // Clean up the path
    return filepath.Clean(path), nil
}