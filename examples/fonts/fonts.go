// fonts demonstrates how to render text using TTF fonts.
package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"time"

	"github.com/mewkiz/pkg/goutil"
	"github.com/mewmew/sdl/font"
	"github.com/mewmew/sdl/win"
	"github.com/mewmew/we"
)

func main() {
	err := fonts()
	if err != nil {
		log.Fatalln(err)
	}
}

// fonts demonstrates how to render text using TTF fonts.
func fonts() (err error) {
	// Open the window.
	err = win.Open(640, 480, win.Resizeable)
	if err != nil {
		return err
	}
	defer win.Close()

	// Load font and image resources.
	err = loadResources()
	if err != nil {
		return err
	}
	defer freeResources()

	// Set the color of the text font to white (the default is black).
	textFont.SetColor(color.White)

	// Render an image of the text "TTF fonts" using the text font.
	textImg, err := textFont.Render("TTF fonts")
	if err != nil {
		return err
	}
	defer textImg.Free()

	// start and frames will be used to calculate the average FPS of the
	// application.
	start := time.Now()
	frames := 0.0

	// Update and event loop.
	for {
		// Render the background image and the rendered text onto the window.
		err = render(textImg)
		if err != nil {
			return err
		}

		// Render the average FPS onto the window.
		err = renderFPS(start, frames)
		if err != nil {
			return err
		}

		// Display window updates on screen.
		err = win.Update()
		if err != nil {
			return err
		}
		frames++

		// Poll events until the event queue is empty.
		for e := win.PollEvent(); e != nil; e = win.PollEvent() {
			fmt.Printf("%T event: %v\n", e, e)
			switch e.(type) {
			case we.Close:
				// Close the application.
				return nil
			}
		}

		// Cap refresh rate at 500 FPS.
		time.Sleep(time.Second / 500)
	}
}

// render renders the background image and the rendered text onto the window.
func render(textImg *win.Image) (err error) {
	// Draw the entire background image onto the screen starting at the top left
	// point (0, 0).
	dp := image.ZP
	err = win.Draw(dp, bgImg)
	if err != nil {
		return err
	}

	// Draw the entire text image onto the screen starting at the destination
	// point (420, 12).
	dp = image.Pt(420, 12)
	err = win.Draw(dp, textImg)
	if err != nil {
		return err
	}

	return nil
}

// renderFPS renders the average FPS to the upper left corner of the window.
func renderFPS(start time.Time, frames float64) (err error) {
	fps := getFPS(start, frames)
	fpsImg, err := fpsFont.Render(fps)
	if err != nil {
		return err
	}
	// Draw the entire fps image onto the screen starting at the destination
	// point (8, 4).
	dp := image.Pt(8, 4)
	win.Draw(dp, fpsImg)
	return nil
}

// getFPS returns the average FPS as a string, based on the provided start time
// and frame count.
func getFPS(start time.Time, frames float64) (text string) {
	// Total runtime in seconds.
	seconds := float64(time.Since(start)) / float64(time.Second)
	// Average FPS.
	fps := frames / seconds
	return fmt.Sprintf("FPS: %.2f", fps)
}

// Background image.
var bgImg *win.Image

// TTF fonts.
var fpsFont, textFont *font.Font

// loadResources loads font and image resources.
func loadResources() (err error) {
	dataDir, err := goutil.SrcDir("github.com/mewmew/sdl/examples/fonts/data")
	if err != nil {
		return err
	}

	// Load background image.
	bgImg, err = win.LoadImage(dataDir + "/bg.png")
	if err != nil {
		return err
	}

	// Load FPS font.
	fpsFont, err = font.Load(dataDir+"/DejaVuSansMono.ttf", 14)
	if err != nil {
		return err
	}
	fpsFont.SetColor(color.White)
	// Use the blended rendering mode for the font.
	fpsFont.SetMode(font.Blended)

	// Load text font.
	textFont, err = font.Load(dataDir+"/Exocet.ttf", 32)
	if err != nil {
		return err
	}
	textFont.SetMode(font.Blended)

	return nil
}

// freeResources frees the memory of font and image resources.
func freeResources() {
	textFont.Free()
	fpsFont.Free()
	bgImg.Free()
}
