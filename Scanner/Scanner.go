package Scanner

import (
	"log"
	"net"
	"strconv"
	"time"
)

// ScanPort Takes protocol, hostname and port to scan port, returns boolean
func ScanPort(protocol, hostname string, port int) bool {
	address := hostname + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout(protocol, address, 3*time.Second)

	if err != nil {
		return false
	}

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(conn)
	return true
}
