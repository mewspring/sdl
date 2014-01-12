// simple demonstrates how to draw surfaces using the Draw and DrawRect methods.
// It also gives an example of a basic event loop.
package main

import (
	"fmt"
	"image"
	"log"
	"time"

	"github.com/mewmew/sdl"
	"github.com/mewmew/we"
)

func main() {
	err := simple()
	if err != nil {
		log.Fatalln(err)
	}
}

// simple demonstrates how to draw surfaces using Draw and DrawRect methods. It
// also gives an example of a basic event loop.
func simple() (err error) {
	// Initialize SDL.
	err = sdl.Init(sdl.InitVideo)
	if err != nil {
		return err
	}
	defer sdl.Quit()

	// Open window.
	win, err := sdl.OpenWindow(640, 480)
	if err != nil {
		return err
	}
	defer win.Close()
	dst, err := win.Surface()
	if err != nil {
		return err
	}

	// Load background surface.
	bg, err := sdl.LoadSurface("data/bg.png")
	if err != nil {
		return err
	}
	defer bg.Free()

	// Draw the entire background surface onto the screen starting at the top
	// left point (0, 0).
	dp := image.ZP
	bg.Draw(dst, dp)

	// Load foreground surface.
	fg, err := sdl.LoadSurface("data/fg.png")
	if err != nil {
		return err
	}
	defer fg.Free()

	// Fill the destination rectangle ((10, 10), (200, 200)) of the screen with
	// corresponding pixels from the foreground surface starting at the source
	// point (70, 70).
	dr := image.Rect(10, 10, 200, 200)
	sp := image.Pt(70, 70)
	fg.DrawRect(dst, dr, sp)

	// start and frames will be used to calculate the average FPS of the
	// application.
	start := time.Now()
	frames := 0

	// displayFPS calculates and displays the average FPS.
	displayFPS := func() {
		seconds := float64(time.Since(start)) / float64(time.Second)
		fps := float64(frames) / seconds
		fmt.Println()
		fmt.Println("=== [ statistics ] =============================================================")
		fmt.Println()
		fmt.Printf("   Total runtime: %.2f seconds.\n", seconds)
		fmt.Printf("   Frame count:   %d frames\n", frames)
		fmt.Printf("   Average FPS:   %.2f frames/second\n", fps)
	}
	defer displayFPS()

	// Draw and event loop.
	for {
		// Update screen.
		err = win.Update()
		if err != nil {
			return err
		}
		frames++

		// Poll events until the event queue is empty.
		for {
			e := sdl.PollEvent()
			if e == nil {
				// The event queue is empty, break loop.
				break
			}
			fmt.Printf("%T event: %v\n", e, e)
			switch e.(type) {
			case we.Close:
				// Close application.
				return nil
			}
		}

		// Cap refresh rate at 500 FPS.
		time.Sleep(time.Second / 500)
	}
}
