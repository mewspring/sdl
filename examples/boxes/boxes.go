// boxes demonstrates how to render text within the confines of a fixed width
// box.
package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"time"

	"github.com/mewkiz/pkg/goutil"
	"github.com/mewmew/sdl/font"
	"github.com/mewmew/sdl/font/fontutil"
	"github.com/mewmew/sdl/window"
	"github.com/mewmew/wandi"
	"github.com/mewmew/we"
)

func main() {
	err := boxes()
	if err != nil {
		log.Fatalln(err)
	}
}

const s = `The Go programming language is an open source project to make programmers more productive.

Go is expressive, concise, clean, and efficient. Its concurrency mechanisms make it easy to write programs that get the most out of multicore and networked machines, while its novel type system enables flexible and modular program construction. Go compiles quickly to machine code yet has the convenience of garbage collection and the power of run-time reflection. It's a fast, statically typed, compiled language that feels like a dynamically typed, interpreted language.`

// boxes demonstrates how to render text within the confines of a fixed width
// box.
func boxes() (err error) {
	// Open the window.
	win, err := window.Open(640, 480, window.Resizeable)
	if err != nil {
		return err
	}
	defer win.Close()

	// Locate data directory.
	dataDir, err := goutil.SrcDir("github.com/mewmew/sdl/examples/boxes/data")
	if err != nil {
		return err
	}

	// Load the background image.
	bg, err := window.LoadImage(dataDir + "/bg.png")
	if err != nil {
		return err
	}
	defer bg.Free()

	// Load the gopher image.
	// ref: http://img.stanleylieber.com/?tags=golang
	goper, err := window.LoadImage(dataDir + "/gopher.png")
	if err != nil {
		return err
	}
	defer goper.Free()

	// Load the font.
	f, err := font.Load(dataDir+"/DejaVuSerif.ttf", 16)
	if err != nil {
		return err
	}
	defer f.Free()

	// Render the text within a bounding box of width 300. It wraps the text so
	// that no single line overflows the pixel width of the box.
	width := 300
	box := fontutil.NewBox(f, width)
	f.SetColor(color.White)
	text, err := box.Render(s)
	if err != nil {
		return err
	}

	// Render the background, gopher and text onto the window.
	err = render(win, bg, goper, text)
	if err != nil {
		return err
	}

	// Update and events loop.
	for {
		// Poll events from the event queue until its empty.
		for e := win.PollEvent(); e != nil; e = win.PollEvent() {
			fmt.Printf("%T event: %v\n", e, e)
			switch e.(type) {
			case we.Close:
				return nil
			case we.Resize:
				// Render the background, gopher and text onto the window.
				err = render(win, bg, goper, text)
				if err != nil {
					return err
				}
			}
		}

		// Display window updates on the screen.
		err = win.Update()
		if err != nil {
			return err
		}

		// Cap the refresh rate at 500 FPS.
		time.Sleep(time.Second / 500)
	}
}

// render renders the background image and the text onto the window.
func render(win wandi.Window, bg, gopher, text *window.Image) (err error) {
	// Render the background onto the window starting at the top left corner
	// (0, 0).
	dp := image.ZP
	err = win.Draw(dp, bg)
	if err != nil {
		return err
	}

	// Render the gopher onto the window starting at the destination point
	//(400, 50).
	dp = image.Pt(400, 50)
	err = win.Draw(dp, gopher)
	if err != nil {
		return err
	}

	// Render the box of text onto the window starting at the destination point
	// (30, 30).
	dp = image.Pt(30, 30)
	err = win.Draw(dp, text)
	if err != nil {
		return err
	}

	return nil
}
