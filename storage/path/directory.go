package path

import (
	"os"
	"path/filepath"
)

func EnsureDir(path string) error {
	dir := filepath.Dir(path)

	if dir == "." {
		dir = path
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}

	return nil
}
