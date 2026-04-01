package path

import (
	"io/fs"
	"path/filepath"
	"regexp"
)

// Files walks dir recursively and returns file paths that match pattern.
//
// If pattern is non-empty it is treated as a regular expression and matched
// against the walked path. If pattern is empty, all files are returned. If
// exclude is non-empty, that path is skipped from the returned results.
func Files(dir, pattern, exclude string) ([]string, error) {
	re, err := regex(pattern)
	if err != nil {
		return nil, err
	}

	excluded, err := abs(exclude)
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

		if excluded != "" {
			normalized, err := filepath.Abs(path)
			if err != nil {
				return err
			}

			if normalized == excluded {
				return nil
			}
		}

		if re == nil || re.MatchString(path) {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func abs(path string) (string, error) {
	if len(path) == 0 {
		return "", nil
	}

	return filepath.Abs(path)
}

func regex(pattern string) (*regexp.Regexp, error) {
	if len(pattern) > 0 {
		return regexp.Compile(pattern)
	}
	return nil, nil
}
