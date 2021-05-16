package ctxio

import (
	"context"
	"io"
)

// ContextWriter io.WriteCloser which can be stopped using context.
type ContextWriter struct {
	io.Writer
	*contextCloser
}

// NewContextWriter creates new ContextWriter connected with context.Context.
func NewContextWriter(ctx context.Context, w io.WriteCloser) *ContextWriter {
	return &ContextWriter{
		Writer:        w,
		contextCloser: newContextCloser(ctx, w),
	}
}
