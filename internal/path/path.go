package path

import (
	"io/fs"
	"path/filepath"
	"regexp"
)

// Files that match a pattern in dir.
func Files(dir, pattern string) ([]string, error) {
	re, err := regex(pattern)
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

		if re == nil || re.MatchString(path) {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func regex(pattern string) (*regexp.Regexp, error) {
	if len(pattern) > 0 {
		return regexp.Compile(pattern)
	}

	return nil, nil
}
