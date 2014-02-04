// simple demonstrates how to draw images using the Draw and DrawRect methods.
// It also gives an example of a basic event loop.
package main

import (
	"fmt"
	"image"
	"log"
	"time"

	"github.com/mewkiz/pkg/goutil"
	"github.com/mewmew/sdl/win"
	"github.com/mewmew/we"
)

func main() {
	err := simple()
	if err != nil {
		log.Fatalln(err)
	}
}

// simple demonstrates how to draw images using the Draw and DrawRect methods.
// It also gives an example of a basic event loop.
func simple() (err error) {
	// Open the window.
	err = win.Open(640, 480, win.Resizeable)
	if err != nil {
		return err
	}
	defer win.Close()

	// Load image resources.
	err = loadResources()
	if err != nil {
		return err
	}
	defer freeResources()

	// start and frames will be used to calculate the average FPS of the
	// application.
	start := time.Now()
	frames := 0

	// Render the images onto the window.
	err = render()
	if err != nil {
		return err
	}

	// Update and event loop.
	for {
		// Display window updates on screen.
		win.Update()
		frames++

		// Poll events until the event queue is empty.
		for e := win.PollEvent(); e != nil; e = win.PollEvent() {
			fmt.Printf("%T event: %v\n", e, e)
			switch e.(type) {
			case we.Close:
				displayFPS(start, frames)
				// Close the application.
				return nil
			case we.Resize:
				// Rerender the images onto the window after resize events.
				err = render()
				if err != nil {
					return err
				}
			}
		}

		// Cap refresh rate at 60 FPS.
		time.Sleep(time.Second / 60)
	}
}

// render renders the background and foreground images onto the window.
func render() (err error) {
	// Draw the entire background image onto the screen starting at the top left
	// point (0, 0).
	dp := image.ZP
	err = win.Draw(dp, bgImg)
	if err != nil {
		return err
	}

	// Fill the destination rectangle ((10, 10), (200, 200)) of the screen with
	// corresponding pixels from the foreground iamge starting at the source
	// point (70, 70).
	dr := image.Rect(10, 10, 200, 200)
	sp := image.Pt(70, 70)
	err = win.DrawRect(dr, fgImg, sp)
	if err != nil {
		return err
	}

	return nil
}

// Background and foreground images.
var bgImg, fgImg *win.Image

// loadResources loads the background and foreground images.
func loadResources() (err error) {
	dataDir, err := goutil.SrcDir("github.com/mewmew/sdl/examples/simple/data")
	if err != nil {
		return err
	}
	// Load background image.
	bgImg, err = win.LoadImage(dataDir + "/bg.png")
	if err != nil {
		return err
	}
	// Load foreground image.
	fgImg, err = win.LoadImage(dataDir + "/fg.png")
	if err != nil {
		return err
	}
	return nil
}

// freeResources frees the memory of the background and foreground images.
func freeResources() {
	fgImg.Free()
	bgImg.Free()
}

// displayFPS calculates and displays the average FPS based on the provided
// frame count.
func displayFPS(start time.Time, frames int) {
	seconds := float64(time.Since(start)) / float64(time.Second)
	fps := float64(frames) / seconds
	fmt.Println()
	fmt.Println("=== [ statistics ] =============================================================")
	fmt.Println()
	fmt.Printf("   Total runtime: %.2f seconds.\n", seconds)
	fmt.Printf("   Frame count:   %d frames\n", frames)
	fmt.Printf("   Average FPS:   %.2f frames/second\n", fps)
}
