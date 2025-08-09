package flag

import (
	"flag"

	"github.com/alexfalkowski/gocovmerge/v2/internal/path"
)

// Files from dir or flag.
func Files(dir, pattern string) ([]string, error) {
	if len(dir) > 0 {
		return path.Files(dir, pattern)
	}
	return flag.Args(), nil
}

// Parse is an alias for flag.Parse.
func Parse() {
	flag.Parse()
}

// StringVar is an alias for flag.StringVar.
func StringVar(p *string, name string, value string, usage string) {
	flag.StringVar(p, name, value, usage)
}
