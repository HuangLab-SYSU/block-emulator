package networks

import (
	"bytes"
	"io"
	"log"
	"net"
	"sync"
)

var connMaplock sync.Mutex
var connectionPool = make(map[string]net.Conn, 0)

func TcpDial(context []byte, addr string) {
	connMaplock.Lock()
	defer connMaplock.Unlock()

	var err error
	var conn net.Conn // Define conn here
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

	_, err = conn.Write(append(context, '\n'))
	if err != nil {
		log.Println("Write error", err)
		return
	}
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

	buffer := make([]byte, 1024)
	var messageBuffer bytes.Buffer

	for {
		n, err := conn.Read(buffer)
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
