package ctxio

import (
	"context"
	"io"
)

// ContextReader io.ReadCloser which can be stopped using context.
type ContextReader struct {
	io.Reader
	*contextCloser
}

// NewContextReader creates new ContextReader connected with context.Context.
func NewContextReader(ctx context.Context, r io.ReadCloser) *ContextReader {
	return &ContextReader{
		Reader:        r,
		contextCloser: newContextCloser(ctx, r),
	}
}
