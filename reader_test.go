package ctxio

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/suite"
)

type readResult struct {
	data []byte
	err  error
}

type ReaderTestSuite struct {
	suite.Suite
}

func (s *ReaderTestSuite) TestReadOK() {
	pr, pw := io.Pipe()
	r := NewContextReader(context.Background(), pr)

	results := make(chan *readResult)
	go func() {
		pw.Write([]byte("a"))
	}()
	go func() {
		res := &readResult{data: make([]byte, 1)}
		_, res.err = r.Read(res.data)
		results <- res
	}()

	// If Read function was not stopped block here.
	res := <-results
	s.NoError(res.err)
	s.EqualValues([]byte("a"), res.data)

	s.NoError(r.Close())
}

func (s *ReaderTestSuite) TestReadClose() {
	pr, _ := io.Pipe()
	r := NewContextReader(context.Background(), pr)

	results := make(chan *readResult)
	go func() {
		res := &readResult{data: make([]byte, 1)}
		_, res.err = r.Read(res.data)
		results <- res
	}()

	s.NoError(r.Close())

	// If Read function was not stopped block here.
	res := <-results
	s.Equal(io.ErrClosedPipe, res.err)
}

func (s *ReaderTestSuite) TestReadCancel() {
	ctx, cancel := context.WithCancel(context.Background())

	pr, _ := io.Pipe()
	r := NewContextReader(ctx, pr)

	results := make(chan *readResult)
	go func() {
		res := &readResult{data: make([]byte, 1)}
		_, res.err = r.Read(res.data)
		results <- res
	}()

	cancel()

	// If Read function was not stopped block here.
	res := <-results
	s.Equal(io.ErrClosedPipe, res.err)

	s.NoError(r.Close())
}

func TestReaderSuite(t *testing.T) {
	suite.Run(t, new(ReaderTestSuite))
}
