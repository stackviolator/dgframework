package main

import (
	"bufio"
	"fmt"
	hb "goHeartBleed/Heartbeat"
	scanner "goHeartBleed/Scanner"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/common-nighthawk/go-figure"
)

// Global vars for super cool colors
var colorGreen = "\033[32m"
var colorReset = "\033[0m"
var colorRed = "\033[31m"

func getCommand() bool {
	var command []string
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n[" + string(colorGreen) + "*" + colorReset + "]gsf > ")

	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalln(err)
	}

	command = strings.Fields(input)

	running := handleCommand(command)

	return running
}

// Gets and handles input for commands
func handleCommand(cmd []string) bool {
	var hostname string
	var portNumbers []string
	var portIntegers []int
	numRoutines := 1000

	// switch statement to handle commands
	switch cmd[0] {
	case "clear":
		fmt.Println("\033[H\033[2J")
	case "help":
		printHelp()
	case "scan":
		routinesInArray, routinesIndex := checkContains(cmd, "-T")
		verbose, _ := checkContains(cmd, "-v")
		hostInArray, hostIndex := checkContains(cmd, "-h")
		portInArray, portIndex := checkContains(cmd, "-p")
		if routinesInArray {
			if cmd[routinesIndex+1] == "MAX" {
				numRoutines = 10000
			}
			numRoutines, _ = strconv.Atoi(cmd[routinesIndex+1])
			if numRoutines > 10000 {
				numRoutines = 10000
			}
		}
		if hostInArray && portInArray {
			hostname = cmd[hostIndex+1]
			if cmd[portIndex+1] == "ALL" {
				portNumbers = append(portNumbers, "1")
				portNumbers = append(portNumbers, "65536")
			}
			if strings.Contains(cmd[portIndex+1], ",") {
				portArray := strings.Split(cmd[portIndex+1], ",")
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

		for j, number := range portNumbers {
			intToAdd, _ := strconv.Atoi(number)
			portIntegers = append(portIntegers, intToAdd)
			j = j
		}

		ports := make(chan int, numRoutines)
		results := make(chan int)
		var openports []int
		var closedports []int

		if len(portIntegers) < 2 {
			portIntegers = append(portIntegers, portIntegers[0])
		}

		// If there is more than one port specified
		sec := time.Now().UnixNano()

		for i := 0; i < cap(ports); i++ {
			go runScan(hostname, ports, verbose, results)
		}

		go func() {
			for i := portIntegers[0]; i <= portIntegers[1]; i++ {
				ports <- i
			}
		}()

		for i := portIntegers[0]; i <= portIntegers[1]; i++ {
			port := <-results
			if port != 0 {
				openports = append(openports, port)
			} else {
				closedports = append(closedports, port)
			}
		}

		close(ports)
		close(results)

		sort.Ints(openports)
		for _, port := range openports {
			fmt.Println("Open port found at "+colorGreen+hostname+":"+strconv.Itoa(port), colorReset)
		}

		difference := (time.Now().UnixNano() - sec) / 1000000000

		fmt.Println("Scanned", portIntegers[1]-portIntegers[0]+1, "ports in", difference, "seconds")

		return true
	case "quit":
		os.Exit(0)
	default:
		printHelp()
	}
	return true
}

func printHelp() {

	fmt.Println("Usage:\n\t[Command] [Options]")
	fmt.Println("----HELP----:\n\tWill display this message, have fun, go crazy")
	fmt.Println("----SCAN----:\n\t-h: Used to specify the host name (full domain or IP)\n\t-p: Used to specify the port(s), a range can be specified with two ports separated by commas (-p 1,100), or ALL for all ports\n\t-T: Used to specifc the amount of goroutines used. Use MAX for 10000+ routines (speeds up scan)\n\t-v: Verbose mode")
	fmt.Println("----QUIT----\n\tWill quit the program")
	fmt.Println("----CLEAR----\n\tI feel like this is very self-explanatory")
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

func runScan(hostname string, ports chan int, verbose bool, results chan int) {
	for p := range ports {
		open := scanner.ScanPort("tcp", hostname, p)
		if open {
			if verbose {
				fmt.Println("Port", strconv.Itoa(p)+colorGreen, "Open", colorReset)
			}
			results <- p
			continue
		} else {
			results <- 0
			if verbose {
				fmt.Println("Port", strconv.Itoa(p)+colorRed, "Closed", colorReset)
			}
		}
	}
}

// Welcome message
func displayWelcomeMessage() {

	fmt.Print(colorRed)
	// git repo has lots of fonts to choose from
	banner := figure.NewFigure("goScan Framework", "larry3d", true)
	banner.Print()

	fmt.Println(colorReset + "Thank you for using my tool it make me happy thinking people are looking at this :) <3\nContact me via email: jpm7050@psu.edu or joshmerrill255@gmail.com <3")
}

func main() {
	argc := len(os.Args[1:])

	// for just sweet cli action
	if argc > 0 {
		displayWelcomeMessage()
		handleCommand(os.Args[1:])
		os.Exit(0)
	}

	displayWelcomeMessage()
	// for full cli interactive mode
	running := getCommand()
	for running {
		running = getCommand()
	}

	hb.Send_hb()
}
