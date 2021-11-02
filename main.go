package main

import (
	"bufio"
	"fmt"
	scanner "goHeartBleed/Scanner"
	"log"
	"os"
	"strconv"
	"strings"
)

// Global vars for super cool colors
var colorGreen = "\033[32m"
var colorReset = "\033[0m"
var colorRed = "\033[31m"

// Gets and handles input for commands
func getCommand() bool {
	reader := bufio.NewReader(os.Stdin)
	var hostname string
	var portNumbers []string
	var portIntegers []int

	fmt.Println("\n[" + string(colorGreen) + "*" + colorReset + "] Enter a command:")

	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}
	cmd := strings.Fields(input)

	switch cmd[0] {
	case "help":
		printHelp()
	case "scan":
		hostInArray, hostIndex := checkContains(cmd, "-h")
		portInArray, portIndex := checkContains(cmd, "-p")
		if hostInArray && portInArray {
			hostname = cmd[hostIndex+1]
			if cmd[portIndex+1] == "ALL" {
				portNumbers = append(portNumbers, "1")
				portNumbers = append(portNumbers, "65536")
			}
			if strings.Contains(cmd[portIndex + 1], ",") {
				portArray := strings.Split(cmd[portIndex + 1], ",")
				for i := range portArray {
					portNumbers = append(portNumbers, portArray[i])
				}
			} else {
				portNumbers = append(portNumbers, cmd[portIndex+1])
			}
		} else if !hostInArray || !portInArray {
			fmt.Println("Invalid syntax, check the 'help' command")
			return true
		}
		//fmt.Println(portNumbers)

		for j, number := range portNumbers {
			intToAdd, _ := strconv.Atoi(number)
			portIntegers = append(portIntegers, intToAdd)
			j = j
		}

		if len(portIntegers) >= 2 {
			for i := portIntegers[0]; i <= portIntegers[1]; i++ {
				runScan(hostname, strconv.Itoa(i))
			}
		} else {
			runScan(hostname, strconv.Itoa(portIntegers[0]))
		}

		return true
	case "quit":
		return false
	default:
		printHelp()
	}
	return true
}

func printHelp() {

	fmt.Println("Usage:\n\t[Command] [Options]")
	fmt.Println("----HELP----:\n\tWill display this message, have fun, go crazy")
	fmt.Println("----SCAN----:\n\t-h: Used to specify the host name (full domain or IP)\n\t-p: Used to specify the port(s), a range can be specified with two ports separated by commas (-p 1,100), or ALL for all ports ")
	fmt.Println("----QUIT----\n\tWill quit the program")
}

// loops through the user inputted command and checks for a flag, returning a boolean and index in the array
func checkContains(arr []string, str string) (bool, int) {
	for k, a := range arr {
		if a == str {
			return true, k
		}
	}
	return false, -1
}

func main() {

	displayWelcomeMessage()
	running := getCommand()
	for running {
		running = getCommand()
	}
}

func runScan(hostname string, port string) bool {
	portNumber, _ := strconv.Atoi(port)
	fmt.Println("Scanning host...", hostname + ":" + strconv.Itoa(portNumber))
	open := scanner.ScanPort("tcp", hostname, portNumber)

	if open {
		fmt.Println("Open port found at "+colorGreen+hostname+":"+port, colorReset)
		return true
	} else {
		// fmt.Println("Port", port + colorRed, "Closed", colorReset)  TODO if verbose mode is added print this
		return false
	}
}

// Welcome message
func displayWelcomeMessage() {
	fmt.Println(colorRed, "   ____               ___ ___                           __       __________ .__                      .___    _________  .____     .___  \n  / ___\\  ____       /   |   \\   ____  _____  _______ _/  |_     \\______   \\|  |    ____   ____    __| _/    \\_   ___ \\ |    |    |   | \n / /_/  >/  _ \\     /    ~    \\_/ __ \\ \\__  \\ \\_  __ \\\\   __\\     |    |  _/|  |  _/ __ \\_/ __ \\  / __ |     /    \\  \\/ |    |    |   | \n \\___  /(  <_> )    \\    Y    /\\  ___/  / __ \\_|  | \\/ |  |       |    |   \\|  |__\\  ___/\\  ___/ / /_/ |     \\     \\____|    |___ |   | \n/_____/  \\____/      \\___|_  /  \\___  >(____  /|__|    |__|       |______  /|____/ \\___  >\\___  >\\____ |      \\______  /|_______ \\|___| \n                           \\/       \\/      \\/                           \\/            \\/     \\/      \\/             \\/         \\/     ")
	fmt.Println(colorReset + "Thank you for using my tool it make me happy thinking people are looking at this :) <3\nContact me via email: jpm7050@psu.edu")
}

// this is a test
