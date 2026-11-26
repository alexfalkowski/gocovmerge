# Accepted design for coverage merging

Part of the [accepted-design index](../accepted-design.md).

## Source revision and file identity

Callers must supply profiles generated from the same source revision because
coverage block coordinates depend on the source. The merge groups entries by
the exact `Profile.FileName` string and orders the result by filename. Different
path spellings for the same source file remain separate entries. Keep this
identity rule independent of the CLI's path comparison for output exclusion.

[Usage constraints](../../../README.md) ·
[Profile merge implementation](../../../internal/cover/cover.go) ·
[Profile merge tests](../../../internal/cover/cover_test.go)

## Modes and block compatibility

A merged profile uses one mode across its files. `set` combines counts with
bitwise OR. `count` and `atomic` sum counts. Reject unknown or mixed modes among
profile entries, and fail if there are no entries to write.

For a given file, blocks can combine when their start and end coordinates and
`NumStmt` values all match. Non-overlapping blocks can be inserted in source
order. Conflicting coordinates, overlapping blocks, or mismatched statement
counts fail the merge because they cannot describe a single consistent set of
coverage blocks.

[Profile merge implementation](../../../internal/cover/cover.go) ·
[Profile merge tests](../../../internal/cover/cover_test.go) ·
[CLI tests](../../../main_test.go)
