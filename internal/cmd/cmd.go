package cmd

import (
	"fmt"
	"io"

	"github.com/alexfalkowski/gocovmerge/v2/internal/cover"
)

// Run parses and merges one or more Go coverage profiles and writes the merged
// profile to out.
//
// files should contain paths to files produced by `go test -coverprofile`.
// Profiles are merged by filename; blocks are combined according to the profile
// mode (set, count, atomic).
//
// It returns an error if any input profile cannot be parsed, if profiles cannot
// be merged (for example, due to mismatched modes or overlapping blocks), or if
// the merged output cannot be written.
func Run(files []string, out io.Writer) error {
	var merged []*cover.Profile

	for _, file := range files {
		profiles, err := cover.ParseProfiles(file)
		if err != nil {
			return fmt.Errorf("failed to parse profiles: %w", err)
		}

		for _, p := range profiles {
			merged, err = cover.AddProfile(merged, p)
			if err != nil {
				return fmt.Errorf("failed to add profile: %w", err)
			}
		}
	}

	if err := cover.WriteProfiles(merged, out); err != nil {
		return fmt.Errorf("failed to write profile: %w", err)
	}

	return nil
}
