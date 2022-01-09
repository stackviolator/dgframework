// ripped from https://stackoverflow.com/questions/40159137/golang-reading-from-stdin-how-to-detect-special-keys-enter-backspace-etc lol tytytyty

package goHeartBleed

import (
	"container/list"
	"fmt"
	"strings"

	cursor "github.com/atomicgo/cursor"
	term "github.com/nsf/termbox-go"
)

func Listener() string {
	var command string
	displayList := list.New()
	displayList.PushBack("")

	fmt.Println()

	currentNode := displayList.Back()

keyPressListenerLoop:
	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			switch ev.Key {
			case term.KeyEsc:
				break keyPressListenerLoop
			case term.KeyF1:
				fmt.Println("F1 pressed")
			case term.KeyF2:
				fmt.Println("F2 pressed")
			case term.KeyF3:
				fmt.Println("F3 pressed")
			case term.KeyF4:
				fmt.Println("F4 pressed")
			case term.KeyF5:
				fmt.Println("F5 pressed")
			case term.KeyF6:
				fmt.Println("F6 pressed")
			case term.KeyF7:
				fmt.Println("F7 pressed")
			case term.KeyF8:
				fmt.Println("F8 pressed")
			case term.KeyF9:
				fmt.Println("F9 pressed")
			case term.KeyF10:
				fmt.Println("F10 pressed")
			case term.KeyF11:
				fmt.Println("F11 pressed")
			case term.KeyF12:
				fmt.Println("F12 pressed")
			case term.KeyInsert:
				fmt.Println("Insert pressed")
			case term.KeyDelete:
				if currentNode != nil {
					tempNode := currentNode.Prev()
					displayList.Remove(currentNode)
					currentNode = tempNode

					cursor.StartOfLine()
					cursor.ClearLine()
					fmt.Print(listToString(displayList))
				}
			case term.KeyHome:
				fmt.Println("Home pressed")
			case term.KeyEnd:
				fmt.Println("End pressed")
			case term.KeyPgup:
				fmt.Println("Page Up pressed")
			case term.KeyPgdn:
				fmt.Println("Page Down pressed")
			case term.KeyArrowUp:
				fmt.Println("Arrow Up pressed")
			case term.KeyArrowDown:
				fmt.Println("Arrow Down pressed")
			case term.KeyArrowLeft:
				if currentNode != nil && currentNode != displayList.Front() {
					cursor.Move(-1, 0)
					currentNode = currentNode.Prev()
				}
			case term.KeyArrowRight:
				cursor.Left(-1)
			case term.KeySpace:
				if currentNode != nil {
					currentNode = displayList.InsertAfter(" ", currentNode)
				} else {
					displayList.PushBack(" ")
				}
				printCommand(listToString(displayList), displayList.Len())
			case term.KeyBackspace:
				if currentNode != nil {
					tempNode := currentNode.Prev()
					displayList.Remove(currentNode)
					currentNode = tempNode

					cursor.StartOfLine()
					cursor.ClearLine()
					fmt.Print(listToString(displayList))
				}
			case term.KeyEnter:
				command = listToString(displayList)
				fmt.Println("")
				return command
				break keyPressListenerLoop
			case term.KeyTab:
				fmt.Println("Tab pressed")

			default:
				// we only want to read a single character or one key pressed event
				if currentNode != nil {
					currentNode = displayList.InsertAfter(string(ev.Ch), currentNode)
				} else {
					displayList.PushBack(string(ev.Ch))
					currentNode = displayList.Back()
				}
				currentIndex := getElementIndex(displayList, currentNode)
				lastIndex := getElementIndex(displayList, displayList.Back())

				printCommand(listToString(displayList), displayList.Len())
				cursor.Right(lastIndex)
				cursor.Left(lastIndex - currentIndex - 1)
			}

		case term.EventError:
			panic(ev.Err)
		}
	}
	return "help"
}

func printCommand(command string, length int) {
	cursor.Show()
	cursor.Move(-(length), 0)
	cursor.ClearLinesDown(0)
	fmt.Print(command)
}

func listToString(ll *list.List) string {
	var strArr []string

	for element := ll.Front(); element != nil; element = element.Next() {
		str, ok := element.Value.(string)
		if !ok {
			fmt.Println("not a string")
			continue
		}
		strArr = append(strArr, str)
	}

	return strings.Join(strArr, "")
}

func getElementIndex(ll *list.List, element *list.Element) int {
	count := 0
	e := ll.Front()
	if e != nil {
		if e == element {
			return count
		} else {
			e = e.Next()
		}
	}

	return -1
}
