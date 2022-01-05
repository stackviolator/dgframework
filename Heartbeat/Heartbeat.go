package Heartbeat

import (
	"fmt"
	"log"
	"net"
	"bufio"
)

func Send_hb() {
	conn, err := net.Dial("tcp", "google.com:80")

	if err != nil {
		fmt.Println(err)
	}

	b := bufio.NewReader(conn)

	for {
		line, err := b.ReadBytes('\n')
		if err != nil {
			break
		}
		conn.Write(line)
	}

	defer func (conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatalln(err)
		}
	} (conn)
}
