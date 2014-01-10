// simple demonstrates how to draw surfaces using the Draw and DrawRect methods.
package main

import (
	"image"
	"log"

	"github.com/mewmew/sdl"
)

func main() {
	err := simple()
	if err != nil {
		log.Fatalln(err)
	}
}

// simple demonstrates how to draw surfaces using Draw and DrawRect methods.
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

	err = win.Update()
	if err != nil {
		return err
	}

	// Block forever.
	select {}
}
