package flag

import (
	"flag"
	"os"

	"github.com/alexfalkowski/gocovmerge/v2/internal/path"
)

// Parse the args and return the parsed values.
func Parse() (*Values, error) {
	args := os.Args[1:]
	set := flag.NewFlagSet("gocovmerge", flag.ContinueOnError)
	out := set.String("o", "", "output file (if missing stdout)")
	dir := set.String("d", "", "directory of files (if missing paths passed in)")
	pattern := set.String("p", "", "pattern to filter directory (if missing all files)")

	if err := set.Parse(args); err != nil {
		return nil, err
	}

	return &Values{out: *out, dir: *dir, pattern: *pattern, args: set.Args()}, nil
}

// Values returns the parsed values from the command line.
type Values struct {
	out     string
	dir     string
	pattern string
	args    []string
}

// Out returns the output file.
func (v *Values) Out() string {
	return v.out
}

// Files from dir or flag.
func (v *Values) Files() ([]string, error) {
	if len(v.dir) > 0 {
		return path.Files(v.dir, v.pattern)
	}
	return v.args, nil
}
