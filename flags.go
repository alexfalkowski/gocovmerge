package main

import "github.com/alexfalkowski/gocovmerge/v2/internal/flag"

func flags() (string, string, string) {
	var (
		out     string
		dir     string
		pattern string
	)

	flag.StringVar(&out, "o", "", "output file (if missing stdout)")
	flag.StringVar(&dir, "d", "", "directory of files (if missing paths passed in)")
	flag.StringVar(&pattern, "p", "", "pattern to filter directory (if missing all files)")
	flag.Parse()

	return out, dir, pattern
}
