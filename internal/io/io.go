package io

import (
	"bytes"
	"io"
	"os"
)

// Writer is the output sink used by the CLI.
type Writer = io.Writer

// OutputWriter accepts merged profile output and finalizes it when closed.
type OutputWriter interface {
	Writer
	Close() error
}

// Output returns a writer for out.
//
// If out is non-empty, Output buffers the merged profile and only creates (or
// truncates) the file when Close is called. Otherwise it writes directly to
// stdout.
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
