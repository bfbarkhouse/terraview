package cmd

import (
	"os"
	"path/filepath"
	"strings"
)

// ResolvePath converts p to an absolute path, expanding ~ and resolving relative paths.
func ResolvePath(p string) (string, error) {
	if strings.HasPrefix(p, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		p = filepath.Join(home, p[2:])
	}
	return filepath.Abs(p)
}
