package networks

import (
	"context"
	"io"
	"log"
	"net"

	"golang.org/x/time/rate"
)

// Write to conn with bandwidth limit.
type rateLimitedWriter struct {
	writer  io.Writer
	limiter *rate.Limiter
}

// Write method waits for the rate limiter and then writes data
func (w *rateLimitedWriter) Write(p []byte) (int, error) {
	// Calculate the number of bytes to write and wait for the limiter to grant enough tokens
	n := len(p)
	for n > 0 {
		writeByteNum := w.limiter.Burst()
		if writeByteNum > n {
			writeByteNum = n
		}
		if err := w.limiter.WaitN(context.TODO(), writeByteNum); err != nil {
			return 0, err
		}
		n -= writeByteNum
	}
	// Actually write the data
	return w.writer.Write(p)
}

func writeToConn(connMsg []byte, conn net.Conn, limiter *rate.Limiter) {
	// Wrap the connection with rateLimitedWriter
	rateLimitedConn := &rateLimitedWriter{writer: conn, limiter: limiter}

	// writing data to the connection
	_, err := rateLimitedConn.Write(connMsg)
	if err != nil {
		log.Println("Write error", err)
		return
	}
}
