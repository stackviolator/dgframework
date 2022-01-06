package main

import (
	"fmt"
	hb "goHeartBleed/Heartbeat"
	keys "goHeartBleed/Keyboard"
	scanner "goHeartBleed/Scanner"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/common-nighthawk/go-figure"
	term "github.com/nsf/termbox-go"
)

// Global vars for super cool colors
var colorGreen = "\033[32m"
var colorReset = "\033[0m"
var colorRed = "\033[31m"

func getCommand() bool {
	var command []string
	fmt.Print("\n[" + string(colorGreen) + "*" + colorReset + "]DGF -- ")

	input := keys.Listener()

	command = strings.Fields(input)

	if len(command) == 0 {
		command = append(command, "balllsss lol")
	}

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

		// Channel of the size of the the gorountines to be run at once
		// Defines worker pool size
		ports := make(chan int, numRoutines)

		// Channel for the results of a scan
		results := make(chan int)

		// Array for which ports are open
		var openports []int

		// If there is less than one port in the array, append the same port to the end
		// e.g. [80] -> [80,80] this will only scan the one port
		if len(portIntegers) < 2 {
			portIntegers = append(portIntegers, portIntegers[0])
		}

		sec := time.Now().UnixNano()

		// While there are still workers in the pool, run a new goroutine of the scan
		for i := 0; i < cap(ports); i++ {
			go runScan(hostname, ports, verbose, results)
		}

		// Adds a new port number to the list of ports to be scanned
		go func() {
			for i := portIntegers[0]; i <= portIntegers[1]; i++ {
				ports <- i
			}
		}()

		// Sends the result to the var port, if it is 0, then the port is closed
		for i := portIntegers[0]; i <= portIntegers[1]; i++ {
			port := <-results
			if port != 0 {
				openports = append(openports, port)
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
	case "heartbeat":
		hb.Send_hb()
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
	// Since the scan will only run when a new port is added to the queue, the length of the channel is the next port to scan
	// e.g. will run at 1 then 2 then 3... then 50 etc. etc
	for p := range ports {
		open := scanner.ScanPort("tcp", hostname, p)
		if open {
			if verbose {
				fmt.Println("Port", strconv.Itoa(p)+colorGreen, "Open", colorReset)
			}
			// If port is open send results the port #
			results <- p
			continue
		} else {
			// Else send a 0
			results <- 0
		}
	}
}

// Welcome message
func displayWelcomeMessage() {

	fmt.Print(colorRed)
	// git repo has lots of fonts to choose from
	banner := figure.NewFigure("DG Framework", "larry3d", true)
	banner.Print()

	fmt.Println(colorReset + "Thank you for using our tool it makes us happy thinking people are looking at this :) <3\nContact us via email: \nJosh \t- jpm7050@psu.edu or joshmerrill@duck.com <3 \nAndrew \t- adm5859@psu.edu ")
}

func main() {
	argc := len(os.Args[1:])

	// for just sweet cli headless action
	if argc > 0 {
		displayWelcomeMessage()
		handleCommand(os.Args[1:])
		os.Exit(0)
	}

	// for full cli interactive mode (cli app)
	err := term.Init()
	if err != nil {
		panic(err)
	}

	defer term.Close()

	displayWelcomeMessage()
	running := getCommand()
	for running {
		running = getCommand()
	}
}
