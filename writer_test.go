package ctxio

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/suite"
)

type WriterTestSuite struct {
	suite.Suite
}

func (s *WriterTestSuite) TestWriteOK() {
	pr, pw := io.Pipe()
	w := NewContextWriter(context.Background(), pw)

	writeErrors := make(chan error, 1)
	results := make(chan *readResult)
	go func() {
		_, err := w.Write([]byte("a"))
		writeErrors <- err
	}()
	go func() {
		res := &readResult{data: make([]byte, 1)}
		_, res.err = pr.Read(res.data)
		results <- res
	}()

	writeError := <-writeErrors
	s.NoError(writeError)

	// If Read function was not stopped block here.
	res := <-results
	s.NoError(res.err)
	s.EqualValues([]byte("a"), res.data)

	s.NoError(w.Close())
}

func (s *WriterTestSuite) TestWriteClose() {
	_, pw := io.Pipe()
	w := NewContextWriter(context.Background(), pw)

	writeErrors := make(chan error, 1)
	go func() {
		_, err := w.Write([]byte("Hello"))
		writeErrors <- err
	}()

	s.NoError(w.Close())

	// If Write function was not stopped block here.
	writeError := <-writeErrors
	s.Equal(io.ErrClosedPipe, writeError)
}

func (s *WriterTestSuite) TestWriteCancel() {
	ctx, cancel := context.WithCancel(context.Background())

	_, pw := io.Pipe()
	w := NewContextWriter(ctx, pw)

	writeErrors := make(chan error, 1)
	go func() {
		_, err := w.Write([]byte("Hello"))
		writeErrors <- err
	}()

	cancel()

	// If Write function was not stopped block here.
	writeError := <-writeErrors
	s.Equal(io.ErrClosedPipe, writeError)

	s.NoError(w.Close())
}

func TestWriterSuite(t *testing.T) {
	suite.Run(t, new(WriterTestSuite))
}
