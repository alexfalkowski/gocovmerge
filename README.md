[![CircleCI](https://circleci.com/gh/alexfalkowski/gocovmerge.svg?style=svg)](https://circleci.com/gh/alexfalkowski/gocovmerge)

# gocovmerge

gocovmerge takes the results from multiple `go test -coverprofile` runs and merges them into one profile.

## usage

```bash
gocovmerge -help

Usage of gocovmerge:
  -d string
        directory of files (if missing paths passed in)
  -o string
        output file (if missing stdout)
  -p string
        pattern to filter directory (if missing all files)
```

You can only merge profiles that were generated from the same source code.

If there are source lines that overlap or do not merge, the process will exit with an error code.
