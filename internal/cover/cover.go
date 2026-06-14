package cover

import (
	"errors"
	"fmt"
	"io"
	"sort"

	"golang.org/x/tools/cover"
)

// Profile is a type alias for golang.org/x/tools/cover.Profile.
//
// It is re-exported so callers can stay within this package when constructing
// or inspecting merged profile data.
type Profile = cover.Profile

// ProfileBlock is a type alias for golang.org/x/tools/cover.ProfileBlock.
//
// It is re-exported alongside Profile so callers do not need to import
// golang.org/x/tools/cover directly.
type ProfileBlock = cover.ProfileBlock

// ErrInvalidMode is returned when a merge or write operation encounters
// profiles that do not all share the same coverage mode.
var ErrInvalidMode = errors.New("invalid profiles merge with different modes")

// ErrUnsupportedMode is returned when a profile uses an unknown coverage mode.
var ErrUnsupportedMode = errors.New("unsupported covermode")

// ErrEmptyProfiles is returned by WriteProfiles when there are no profiles to
// write.
var ErrEmptyProfiles = errors.New("empty profiles")

// ParseProfiles parses a coverage profile file produced by `go test -coverprofile`.
//
// The returned profiles are sorted by filename, and each profile's blocks are
// sorted by start position using golang.org/x/tools/cover's parser semantics.
func ParseProfiles(fileName string) ([]*Profile, error) {
	return cover.ParseProfiles(fileName)
}

// AddProfile inserts p into profiles, merging it with an existing profile for
// the same filename if present.
//
// profiles must already be sorted by Profile.FileName and contain one coverage
// mode, typically because it was produced by prior AddProfile calls. Existing
// blocks and p.Blocks must also be sorted by block start position, typically by
// ParseProfiles. The returned slice is kept sorted by Profile.FileName.
//
// AddProfile mutates existing profiles when merging counts and inserts p
// directly when adding a new filename; it does not deep-copy profile data. If a
// profile with the same FileName already exists, blocks from p are merged into
// it. Blocks with the same coordinates must also agree on NumStmt; otherwise
// AddProfile returns an error. Overlapping or otherwise incompatible blocks are
// rejected.
func AddProfile(profiles []*Profile, p *Profile) ([]*Profile, error) {
	if err := validateMode(p.Mode); err != nil {
		return nil, err
	}

	if len(profiles) > 0 && profiles[0].Mode != p.Mode {
		return nil, ErrInvalidMode
	}

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

// WriteProfiles writes profiles to out using the standard Go coverage profile
// format.
//
// The output starts with a single `mode: ...` line followed by one line per
// block. Before writing, it validates that every profile uses the same mode. It
// returns ErrEmptyProfiles when profiles is empty.
func WriteProfiles(profiles []*Profile, out io.Writer) error {
	if len(profiles) == 0 {
		return ErrEmptyProfiles
	}

	mode := profiles[0].Mode
	if err := validateMode(mode); err != nil {
		return err
	}

	for _, p := range profiles[1:] {
		if err := validateMode(p.Mode); err != nil {
			return err
		}

		if p.Mode != mode {
			return ErrInvalidMode
		}
	}

	_, err := fmt.Fprintf(out, "mode: %s\n", mode)
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

func mergeProfiles(p, merge *Profile) error {
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

func mergeProfileBlock(p *Profile, pb ProfileBlock, startIndex int) (int, error) {
	if startIndex >= len(p.Blocks) { // no more to merge, append the remaining blocks
		p.Blocks = append(p.Blocks, pb)
		return len(p.Blocks), nil
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

	if i < len(p.Blocks) && p.Blocks[i].StartLine == pb.StartLine && p.Blocks[i].StartCol == pb.StartCol {
		if err := mergeSameBlock(p, i, pb); err != nil {
			return 0, err
		}

		return i + 1, nil
	}

	if err := insertBlock(p, i, pb); err != nil {
		return 0, err
	}

	return i + 1, nil
}

func mergeSameBlock(p *Profile, index int, block ProfileBlock) error {
	if p.Blocks[index].EndLine != block.EndLine || p.Blocks[index].EndCol != block.EndCol {
		return fmt.Errorf("overlap merge: %v %v %v", p.FileName, p.Blocks[index], block)
	}

	if p.Blocks[index].NumStmt != block.NumStmt {
		return fmt.Errorf("inconsistent NumStmt: changed from %d to %d", p.Blocks[index].NumStmt, block.NumStmt)
	}

	switch p.Mode {
	case "set":
		p.Blocks[index].Count |= block.Count
	case "count", "atomic":
		p.Blocks[index].Count += block.Count
	default:
		return fmt.Errorf("%w: %q", ErrUnsupportedMode, p.Mode)
	}

	return nil
}

func validateMode(mode string) error {
	switch mode {
	case "set", "count", "atomic":
		return nil
	default:
		return fmt.Errorf("%w: %q", ErrUnsupportedMode, mode)
	}
}

func insertBlock(p *Profile, index int, block ProfileBlock) error {
	if index > 0 {
		previous := p.Blocks[index-1]
		if previous.EndLine > block.StartLine || (previous.EndLine == block.StartLine && previous.EndCol > block.StartCol) {
			return fmt.Errorf("overlap before: %v %v %v", p.FileName, previous, block)
		}
	}

	if index < len(p.Blocks) {
		next := p.Blocks[index]
		if next.StartLine < block.EndLine || (next.StartLine == block.EndLine && next.StartCol < block.EndCol) {
			return fmt.Errorf("overlap after: %v %v %v", p.FileName, next, block)
		}
	}

	p.Blocks = append(p.Blocks, cover.ProfileBlock{})
	copy(p.Blocks[index+1:], p.Blocks[index:])
	p.Blocks[index] = block

	return nil
}
