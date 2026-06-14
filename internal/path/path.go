package path

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

// Files walks dir recursively and returns file paths that match pattern.
//
// If pattern is non-empty it is treated as a regular expression and matched
// against the walked path. If pattern is empty, all files are returned. If
// exclude is non-empty, that path is skipped from the returned results. The
// exclusion check canonicalizes walked paths and exclude with symlink
// evaluation so relative `-o` values are handled consistently. Missing excluded
// paths are resolved through their canonical parent, but other canonicalization
// errors are returned.
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
			normalized, err := abs(path)
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

	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	canonical, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}

		return filepath.Join(canonicalDir(absolute), filepath.Base(absolute)), nil
	}

	return canonical, nil
}

func canonicalDir(path string) string {
	dir, err := filepath.EvalSymlinks(filepath.Dir(path))
	if err != nil {
		return filepath.Dir(path)
	}

	return dir
}

func regex(pattern string) (*regexp.Regexp, error) {
	if len(pattern) > 0 {
		return regexp.Compile(pattern)
	}
	return nil, nil
}
