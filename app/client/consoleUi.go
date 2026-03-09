package client

import (
	"fmt"
	"sync"
)

const (
	ansiClearLine = "\033[2K"
)

func renderPromptLocked() {
	fmt.Print("> ")
}

func clearCurrentLineLocked() {
	fmt.Print("\r")
	fmt.Print(ansiClearLine)
}

func printLine(outMu *sync.Mutex, line string) {
	outMu.Lock()
	defer outMu.Unlock()

	clearCurrentLineLocked()
	fmt.Println(line)
	renderPromptLocked()
}

func printHelp(outMu *sync.Mutex) {
	outMu.Lock()
	defer outMu.Unlock()

	clearCurrentLineLocked()
	fmt.Println("Commands:")
	fmt.Println("  /help")
	fmt.Println("  /nick NAME")
	fmt.Println("  /rooms")
	fmt.Println("  /create ROOM [PASS]")
	fmt.Println("  /join ROOM [PASS]")
	fmt.Println("  /leave")
	fmt.Println("  /quit")
	renderPromptLocked()
}

func printLocal(outMu *sync.Mutex, text string) {
	printLine(outMu, "[local] "+text)
}

func printPrompt(outMu *sync.Mutex) {
	outMu.Lock()
	defer outMu.Unlock()

	clearCurrentLineLocked()
	renderPromptLocked()
}
