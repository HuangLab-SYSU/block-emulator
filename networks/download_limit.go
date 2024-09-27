package networks

import (
	"context"
	"io"
	"net"

	"golang.org/x/time/rate"
)

// rateLimitedReader wraps an io.Reader with a rate limiter
type rateLimitedReader struct {
	reader  io.Reader
	limiter *rate.Limiter
}

// Read method waits for the rate limiter and then reads data
func (r *rateLimitedReader) Read(p []byte) (int, error) {
	// Calculate the number of bytes to read and wait for the limiter to grant enough tokens
	n := len(p)
	for n > 0 {
		writeByteNum := r.limiter.Burst()
		if writeByteNum > n {
			writeByteNum = n
		}
		if err := r.limiter.WaitN(context.TODO(), n); err != nil {
			return 0, err
		}
		n -= writeByteNum
	}

	// Actually read the data
	return r.reader.Read(p)
}

func NewConnReader(conn net.Conn, rateLimiter *rate.Limiter) *rateLimitedReader {
	// Wrap the connection with rateLimitedReader
	return &rateLimitedReader{reader: conn, limiter: rateLimiter}
}
