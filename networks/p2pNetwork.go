package networks

import (
	"blockEmulator/params"
	"bytes"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"math/rand"

	"golang.org/x/time/rate"
)

var connMaplock sync.Mutex
var connectionPool = make(map[string]net.Conn, 0)

// network params.
var randomDelayGenerator *rand.Rand
var rateLimiterDownload *rate.Limiter
var rateLimiterUpload *rate.Limiter

// Define the latency, jitter and bandwidth here.
// Init tools.
func InitNetworkTools() {
	// avoid wrong params.
	if params.Delay < 0 {
		params.Delay = 0
	}
	if params.JitterRange < 0 {
		params.JitterRange = 0
	}
	if params.Bandwidth < 0 {
		params.Bandwidth = 0x7fffffff
	}

	// generate the random seed.
	randomDelayGenerator = rand.New(rand.NewSource(time.Now().UnixMicro()))
	// Limit the download rate
	rateLimiterDownload = rate.NewLimiter(rate.Limit(params.Bandwidth), params.Bandwidth)
	// Limit the upload rate
	rateLimiterUpload = rate.NewLimiter(rate.Limit(params.Bandwidth), params.Bandwidth)
}

func TcpDial(context []byte, addr string) {
	go func() {
		// simulate the delay
		thisDelay := params.Delay
		if params.JitterRange != 0 {
			thisDelay = randomDelayGenerator.Intn(params.JitterRange) - params.JitterRange/2 + params.Delay
		}
		time.Sleep(time.Millisecond * time.Duration(thisDelay))

		connMaplock.Lock()
		defer connMaplock.Unlock()

		var err error
		var conn net.Conn // Define conn here

		// if this connection is not built, build it.
		if c, ok := connectionPool[addr]; ok {
			if tcpConn, tcpOk := c.(*net.TCPConn); tcpOk {
				if err := tcpConn.SetKeepAlive(true); err != nil {
					delete(connectionPool, addr) // Remove if not alive
					conn, err = net.Dial("tcp", addr)
					if err != nil {
						log.Println("Reconnect error", err)
						return
					}
					connectionPool[addr] = conn
					go ReadFromConn(addr) // Start reading from new connection
				} else {
					conn = c // Use the existing connection
				}
			}
		} else {
			conn, err = net.Dial("tcp", addr)
			if err != nil {
				log.Println("Connect error", err)
				return
			}
			connectionPool[addr] = conn
			go ReadFromConn(addr) // Start reading from new connection
		}

		writeToConn(append(context, '\n'), conn, rateLimiterUpload)
	}()
}

// Broadcast sends a message to multiple receivers, excluding the sender.
func Broadcast(sender string, receivers []string, msg []byte) {
	for _, ip := range receivers {
		if ip == sender {
			continue
		}
		go TcpDial(msg, ip)
	}
}

// CloseAllConnInPool closes all connections in the connection pool.
func CloseAllConnInPool() {
	connMaplock.Lock()
	defer connMaplock.Unlock()

	for _, conn := range connectionPool {
		conn.Close()
	}
	connectionPool = make(map[string]net.Conn) // Reset the pool
}

// ReadFromConn reads data from a connection.
func ReadFromConn(addr string) {
	conn := connectionPool[addr]

	// new a conn reader
	connReader := NewConnReader(conn, rateLimiterDownload)

	buffer := make([]byte, 1024)
	var messageBuffer bytes.Buffer

	for {
		n, err := connReader.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Println("Read error for address", addr, ":", err)
			}
			break
		}

		// add message to buffer
		messageBuffer.Write(buffer[:n])

		// handle the full message
		for {
			message, err := readMessage(&messageBuffer)
			if err == io.ErrShortBuffer {
				// continue to load if buffer is short
				break
			} else if err == nil {
				// log the full message
				log.Println("Received from", addr, ":", message)
			} else {
				// handle other errs
				log.Println("Error processing message for address", addr, ":", err)
				break
			}
		}
	}
}

func readMessage(buffer *bytes.Buffer) (string, error) {
	message, err := buffer.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	return string(message), nil
}
