package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	files := []string{
		filepath.Join("cmd", "x404x", "console.go"),
		filepath.Join("cmd", "x404x", "ui.go"),
		filepath.Join("cmd", "x404x", "commands.go"),
		filepath.Join("cmd", "x404x", "root.go"),
		filepath.Join("cmd", "x404x", "payload_commands.go"),
		filepath.Join("cmd", "x404x", "listeners.go"),
		filepath.Join("cmd", "x404x", "tui.go"),
	}

	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			fmt.Println("Skipping", f, err)
			continue
		}

		b = bytes.ReplaceAll(b, []byte("fmt.Printf("), []byte("fmt.Fprintf(ConsoleOut, "))
		b = bytes.ReplaceAll(b, []byte("fmt.Println("), []byte("fmt.Fprintln(ConsoleOut, "))
		b = bytes.ReplaceAll(b, []byte("fmt.Print("), []byte("fmt.Fprint(ConsoleOut, "))

		err = os.WriteFile(f, b, 0644)
		if err != nil {
			fmt.Println("Error writing", f, err)
		} else {
			fmt.Println("Refactored", f)
		}
	}
}
