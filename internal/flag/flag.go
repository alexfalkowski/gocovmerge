package flag

import (
	"flag"
	"os"

	"github.com/alexfalkowski/gocovmerge/v2/internal/path"
)

// Parse parses command-line flags and returns the resulting Values.
//
// Supported flags:
//   - `-o`: output file path (if empty, stdout is used)
//   - `-d`: directory containing coverage profiles (if empty, positional args are used)
//   - `-p`: regexp pattern to filter files when `-d` is set (if empty, all files are included)
//
// Any remaining positional arguments after flags are treated as coverage profile paths.
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

// Values holds the parsed command-line configuration.
type Values struct {
	out     string
	dir     string
	pattern string
	args    []string
}

// Out returns the output file path provided via `-o`.
//
// If empty, the caller should write to stdout.
func (v *Values) Out() string {
	return v.out
}

// Files returns the coverage profile file paths to merge.
//
// If `-d` was provided, it walks that directory recursively and returns all
// matching files; if `-p` is non-empty it is treated as a regexp filter.
// Otherwise, it returns the remaining positional arguments.
func (v *Values) Files() ([]string, error) {
	if len(v.dir) > 0 {
		return path.Files(v.dir, v.pattern)
	}
	return v.args, nil
}
