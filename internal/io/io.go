package io

import (
	"bytes"
	"io"
	"os"
)

// Writer is the output sink used by the CLI.
//
// It aliases io.Writer so callers can depend on this package without importing
// the standard library's io package directly.
type Writer = io.Writer

// OutputWriter accepts merged profile output and finalizes it when closed.
//
// File-backed implementations buffer writes in memory and only create or
// truncate the destination when Close succeeds. The stdout-backed
// implementation's Close method is a no-op.
type OutputWriter interface {
	Writer
	Close() error
}

// Output returns a writer for out.
//
// If out is non-empty, Output buffers the merged profile and only creates (or
// truncates) the file when Close is called. Otherwise it writes directly to
// stdout and Close becomes a no-op.
func Output(out string, stdout Writer) OutputWriter {
	if len(out) > 0 {
		return &fileOutput{path: out}
	}
	return &stdoutOutput{Writer: stdout}
}

type stdoutOutput struct {
	Writer
}

func (o *stdoutOutput) Close() error {
	return nil
}

type fileOutput struct {
	path   string
	buffer bytes.Buffer
}

func (o *fileOutput) Write(p []byte) (int, error) {
	return o.buffer.Write(p)
}

func (o *fileOutput) Close() error {
	file, err := os.Create(o.path)
	if err != nil {
		return err
	}

	if _, err := o.buffer.WriteTo(file); err != nil {
		_ = file.Close()
		return err
	}

	return file.Close()
}
