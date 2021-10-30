package Scanner

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

// var colorGreen = "\033[32m"
var colorReset = "\033[0m"
var colorRed = "\033[31m"

// ScanPort Takes protocol, hostname and port to scan port, returns boolean
func ScanPort(protocol, hostname string, port int) bool {
	address := hostname + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout(protocol, address, 60*time.Second)

	if err != nil {
		fmt.Print("[" + colorRed + "*" + colorReset + "] Ran into error: ")
		log.Fatalln(err)
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