package path

import (
	"os"
	"path/filepath"
	"strings"
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

func GetFilesInDir(dirPath, pattern string) ([]string, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var matchedFiles []string

	parts := strings.Split(pattern, "*")
	prefix := parts[0]
	suffix := ""

	if len(parts) > 1 {
		suffix = parts[1]
	}

	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			if strings.HasPrefix(fileName, prefix) && strings.HasSuffix(fileName, suffix) {
				matchedFiles = append(matchedFiles, filepath.Join(dirPath, fileName))
			}
		}
	}

	return matchedFiles, nil
}
