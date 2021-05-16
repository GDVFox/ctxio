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
	}
	go ctxCloser.observeCloserContext()

	return ctxCloser
}

func (w *contextCloser) observeCloserContext() {
	<-w.ctx.Done()
	w.errs <- w.Closer.Close()
}

// Close closes underlying Closer.
func (w *contextCloser) Close() error {
	w.once.Do(func() {
		w.stop()
		w.closeErr = <-w.errs
	})
	return w.closeErr
}
