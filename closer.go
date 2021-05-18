package ctxio

import (
	"context"
	"io"
	"sync"
)

// ContextCloser io.Closer which can be stopped using context.
type contextCloser struct {
	io.Closer

	once sync.Once
	ctx  context.Context
	stop context.CancelFunc

	errs     chan error
	closeErr error

	freeCh chan struct{}
}

// NewContextCloser creates new ContextCloser connected with context.Context.
func newContextCloser(ctx context.Context, c io.Closer) *contextCloser {
	closeCtx, closeStop := context.WithCancel(ctx)

	ctxCloser := &contextCloser{
		Closer: c,
		once:   sync.Once{},
		ctx:    closeCtx,
		stop:   closeStop,
		errs:   make(chan error, 1),
		freeCh: make(chan struct{}),
	}
	go ctxCloser.observeCloserContext()

	return ctxCloser
}

func (w *contextCloser) observeCloserContext() {
	select {
	case <-w.ctx.Done():
		w.errs <- w.Closer.Close()
	case <-w.freeCh:
		w.errs <- nil
	}
}

// Free removes connection with context.
// After this call Close or context cancel will not close underlying Closer.
func (w *contextCloser) Free() error {
	w.once.Do(func() {
		close(w.freeCh)
		w.closeErr = <-w.errs

		// Always call context stop to free context.
		// Call stop() after close(w.freeCh) because we do not want to close underlying Closer.
		w.stop()
	})

	return w.closeErr
}

// Close closes underlying Closer.
// After Free() call always returns nil.
func (w *contextCloser) Close() error {
	w.once.Do(func() {
		w.stop()
		w.closeErr = <-w.errs
	})
	return w.closeErr
}
