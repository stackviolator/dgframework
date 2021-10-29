package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

//TODO make the scanport func its own module?

// Global vars for super cool colors
var colorGreen string = "\033[32m"
var colorReset string = "\033[0m"

// Takes protocol, hostname and port to scan port, returns boolean
func scanPort(protocol, hostname string, port int,) bool {
	address := hostname + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout(protocol, address, 60*time.Second)

	if err != nil {
		log.Fatalln(err)
		return false
	}

	defer conn.Close()
	return true
}

// Gets and handles input for commands
//TODO actually process flags in commands lol
func getCommand() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	displayWelcomeMessage()
	fmt.Println("[" + string(colorGreen) + "*" + colorReset + "] Enter a command:")

	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}

	cmd := strings.Fields(input)
	if cmd[0] == "help" {
		fmt.Println("Example usage:\nscan -h google.com -p 443")
	} else if cmd[0] == "scan" {
		return cmd[2], cmd[4]
	}

	return "", ""
}

func main() {
	hostname, port := getCommand()
	portNumber, _ := strconv.Atoi(port)

	fmt.Println(hostname)
	fmt.Println("Scanning port")
	open := scanPort("tcp", hostname, portNumber)

	if open {
		fmt.Println("Open port found at " + colorGreen + hostname + ":" + port, colorReset)
	}

}

func displayWelcomeMessage () {

	fmt.Println("                 ___ ___                        __ __________.__                    .___ _________ .____    .___ \n   ____   ____  /   |   \\   ____ _____ ________/  |\\______   \\  |   ____   ____   __| _/ \\_   ___ \\|    |   |   |\n  / ___\\ /  _ \\/    ~    \\_/ __ \\\\__  \\\\_  __ \\   __\\    |  _/  | _/ __ \\_/ __ \\ / __ |  /    \\  \\/|    |   |   |\n / /_/  >  <_> )    Y    /\\  ___/ / __ \\|  | \\/|  | |    |   \\  |_\\  ___/\\  ___// /_/ |  \\     \\___|    |___|   |\n \\___  / \\____/ \\___|_  /  \\___  >____  /__|   |__| |______  /____/\\___  >\\___  >____ |   \\______  /_______ \\___|\n/_____/               \\/       \\/     \\/                   \\/          \\/     \\/     \\/          \\/        \\/    ")
}