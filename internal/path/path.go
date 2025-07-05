package path

import (
	"io/fs"
	"path/filepath"
	"regexp"
)

// Files that match a pattern in dir.
func Files(dir, pattern string) ([]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	var files []string

	err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if re.MatchString(path) {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}
