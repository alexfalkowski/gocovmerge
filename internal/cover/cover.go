package cover

import (
	"errors"
	"fmt"
	"io"
	"sort"

	"golang.org/x/tools/cover"
)

// Profile is an alias for cover.Profile.
type Profile = cover.Profile

var (
	// ErrInvalidMode for cover.
	ErrInvalidMode = errors.New("invalid profiles merge with different modes")

	// ErrEmptyProfiles for cover.
	ErrEmptyProfiles = errors.New("empty profiles")

	// ParseProfiles is an alias of cover.ParseProfiles.
	ParseProfiles = cover.ParseProfiles
)

// AddProfile for cover.
func AddProfile(profiles []*cover.Profile, p *cover.Profile) ([]*cover.Profile, error) {
	i := sort.Search(len(profiles), func(i int) bool { return profiles[i].FileName >= p.FileName })

	if i < len(profiles) && profiles[i].FileName == p.FileName {
		if err := mergeProfiles(profiles[i], p); err != nil {
			return nil, err
		}
	} else {
		profiles = append(profiles, nil)
		copy(profiles[i+1:], profiles[i:])
		profiles[i] = p
	}

	return profiles, nil
}

// WriteProfiles to out.
func WriteProfiles(profiles []*cover.Profile, out io.Writer) error {
	if len(profiles) == 0 {
		return ErrEmptyProfiles
	}

	_, err := fmt.Fprintf(out, "mode: %s\n", profiles[0].Mode)
	if err != nil {
		return err
	}

	for _, p := range profiles {
		for _, b := range p.Blocks {
			_, err := fmt.Fprintf(out, "%s:%d.%d,%d.%d %d %d\n", p.FileName, b.StartLine, b.StartCol, b.EndLine, b.EndCol, b.NumStmt, b.Count)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func mergeProfiles(p, merge *cover.Profile) error {
	if p.Mode != merge.Mode {
		return ErrInvalidMode
	}

	// Since the blocks are sorted, we can keep track of where the last block
	// was inserted and only look at the blocks after that as targets for merge
	startIndex := 0
	for _, b := range merge.Blocks {
		i, err := mergeProfileBlock(p, b, startIndex)
		if err != nil {
			return err
		}

		startIndex = i
	}

	return nil
}

func mergeProfileBlock(p *cover.Profile, pb cover.ProfileBlock, startIndex int) (int, error) {
	if startIndex >= len(p.Blocks) { // no more to merge
		return startIndex, nil
	}

	sortFunc := func(i int) bool {
		pi := p.Blocks[i+startIndex]

		return pi.StartLine >= pb.StartLine && (pi.StartLine != pb.StartLine || pi.StartCol >= pb.StartCol)
	}

	i := 0
	if !sortFunc(i) {
		i = sort.Search(len(p.Blocks)-startIndex, sortFunc)
	}

	i += startIndex

	//nolint:nestif
	if i < len(p.Blocks) && p.Blocks[i].StartLine == pb.StartLine && p.Blocks[i].StartCol == pb.StartCol {
		if p.Blocks[i].EndLine != pb.EndLine || p.Blocks[i].EndCol != pb.EndCol {
			return 0, fmt.Errorf("overlap merge: %v %v %v", p.FileName, p.Blocks[i], pb)
		}

		switch p.Mode {
		case "set":
			p.Blocks[i].Count |= pb.Count
		case "count", "atomic":
			p.Blocks[i].Count += pb.Count
		default:
			return 0, fmt.Errorf("unsupported covermode: '%s'", p.Mode)
		}
	} else {
		if i > 0 {
			pa := p.Blocks[i-1]
			if pa.EndLine >= pb.EndLine && (pa.EndLine != pb.EndLine || pa.EndCol > pb.EndCol) {
				return 0, fmt.Errorf("overlap before: %v %v %v", p.FileName, pa, pb)
			}
		}

		if i < len(p.Blocks)-1 {
			pa := p.Blocks[i+1]
			if pa.StartLine <= pb.StartLine && (pa.StartLine != pb.StartLine || pa.StartCol < pb.StartCol) {
				return 0, fmt.Errorf("overlap after: %v %v %v", p.FileName, pa, pb)
			}
		}

		p.Blocks = append(p.Blocks, cover.ProfileBlock{})
		copy(p.Blocks[i+1:], p.Blocks[i:])
		p.Blocks[i] = pb
	}

	return i + 1, nil
}
