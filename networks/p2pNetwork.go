// something about broadcast

package networks

import (
	"log"
	"net"
	"sync"
)

var connMaplock sync.Mutex
var connectionPool = make(map[string]net.Conn, 0)

func TcpDial(context []byte, addr string) {
	var conn net.Conn
	connMaplock.Lock()
	defer connMaplock.Unlock()
	if connectionPool[addr] == nil {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			log.Println("connect error", err)
			return
		}
		connectionPool[addr] = conn
	}
	conn = connectionPool[addr]

	_, err := conn.Write(append(context, '\n'))
	if err != nil {
		return
	}
}

func Broadcast(sender string, receivers []string, msg []byte) {
	for _, ip := range receivers {
		if ip == sender {
			continue
		}
		go TcpDial(msg, ip)
	}
}

func CloseAllConnInPool() {
	for _, conn := range connectionPool {
		conn.Close()
	}
}

// todo
// long connect (not close immediately) ...
