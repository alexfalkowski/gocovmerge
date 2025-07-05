package cmd

import (
	"fmt"
	"io"

	"github.com/alexfalkowski/gocovmerge/v2/internal/cover"
)

// Run the command.
func Run(out io.Writer, args []string) error {
	var merged []*cover.Profile

	for _, file := range args {
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
