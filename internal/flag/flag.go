package flag

import (
	"flag"

	"github.com/alexfalkowski/gocovmerge/v2/internal/path"
)

var (
	// Parse is an alias for flag.Parse.
	Parse = flag.Parse

	// StringVar is an alias for flag.StringVar.
	StringVar = flag.StringVar
)

// Files from dir or flag.
func Files(dir, pattern string) ([]string, error) {
	if len(dir) > 0 {
		return path.Files(dir, pattern)
	}

	return flag.Args(), nil
}
