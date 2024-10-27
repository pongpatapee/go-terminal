package main

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/creack/pty"
)

// limit for our command output buffer
const MaxBufferSize = 16

func main() {
	a := app.New()
	w := a.NewWindow("go-terminal")

	ui := widget.NewTextGrid()
	ui.SetText("hello from go terminal")

	// Prepare shell to be openned
	c := exec.Command("/bin/bash")
	// Can run my personal simple shell instead
	// c := exec.Command("/home/dan/personal/coding/projects/go-shell/go-shell")

	// pty: pseudo-terminal
	// A pair of pseudo devices (master and slave, yes very unforunate terminalogies)
	// which estabilishes asynchronous, bidirectional communication
	// Whatever that means

	p, err := pty.Start(c)
	if err != nil {
		fyne.LogError("Failed to open pty", err)
		os.Exit(1)
	}

	defer c.Process.Kill()

	// callback functions for type event

	onTypedKey := func(e *fyne.KeyEvent) {
		if e.Name == fyne.KeyEnter || e.Name == fyne.KeyReturn {
			// Write a carriage return on "enter" key
			_, _ = p.Write([]byte{'\r'})
		}
	}

	onTypedRune := func(r rune) {
		_, _ = p.WriteString(string(r))
	}

	// Set call back function on type events to the
	// Terminal emulator
	w.Canvas().SetOnTypedKey(onTypedKey)
	w.Canvas().SetOnTypedRune(onTypedRune)

	buffer := [][]rune{}
	reader := bufio.NewReader(p)

	// go routine that reads from pty
	go func() {
		line := []rune{}
		buffer = append(buffer, line)

		for {

			r, _, err := reader.ReadRune()
			if err != nil {
				if err == io.EOF {
					return
				}
				os.Exit(0)
			}

			line = append(line, r)
			buffer[len(buffer)-1] = line

			if r == '\n' {
				// pop first item if buffer reaches cap
				if len(buffer) > MaxBufferSize {
					buffer = buffer[1:]
				}

				line = []rune{}
				buffer = append(buffer, line)
			}

		}
	}()

	// go routine that renders UI
	go func() {
		for {

			time.Sleep(100 * time.Millisecond)
			ui.SetText("")
			var lines string

			for _, line := range buffer {
				lines = lines + string(line)
			}

			ui.SetText(string(lines))
		}
	}()

	w.SetContent(
		container.New(
			layout.NewGridWrapLayout(
				fyne.NewSize(900, 325),
			),
			ui,
		),
	)
	w.ShowAndRun()
}
