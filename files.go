package main

import (
	"flag"

	"github.com/alexfalkowski/gocovmerge/v2/internal/path"
)

func files(dir, pattern string) ([]string, error) {
	if len(dir) > 0 {
		f, err := path.Files(dir, pattern)
		if err != nil {
			return nil, err
		}

		return f, nil
	}

	return flag.Args(), nil
}
