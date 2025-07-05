package cmd

import (
	"fmt"
	"io"

	"github.com/alexfalkowski/gocovmerge/v2/internal/cover"
)

// Run the command.
func Run(stdout io.Writer, args []string) error {
	var merged []*cover.Profile

	for _, file := range args {
		profiles, err := cover.ParseProfiles(file)
		if err != nil {
			return fmt.Errorf("cmd: failed to parse profiles: %w", err)
		}

		for _, p := range profiles {
			merged, err = cover.AddProfile(merged, p)
			if err != nil {
				return fmt.Errorf("cmd: failed to add profile: %w", err)
			}
		}
	}

	cover.WriteProfiles(merged, stdout)

	return nil
}
