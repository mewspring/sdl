// Package window handles window creation, drawing and events. It uses a small
// subset of the features provided by the SDL library version 2.0 [1].
//
// [1]: http://www.libsdl.org/
package window

// #cgo pkg-config: sdl2
// #include <SDL2/SDL.h>
import "C"

import (
	"image"
	"image/color"

	"github.com/mewmew/wandi"
)

// Flag is a bitfield of window flags.
type Flag uint32

const (
	// Resizeable states that the window can be resized.
	Resizeable Flag = C.SDL_WINDOW_RESIZABLE
	// Fullscreen states that the window is in full screen mode.
	Fullscreen Flag = C.SDL_WINDOW_FULLSCREEN
)

// A Window represents a graphical window capable of handling draw operations
// and window events. It implements the wandi.Window interface.
type Window struct {
	win *C.SDL_Window
}

// winCount represent the number of active windows.
var winCount int

// Open opens a new window of the specified dimensions and optional window
// flags. By default the window is not resizeable.
//
// Note: The Close method of the window must be called when finished using it.
func Open(width, height int, flags ...Flag) (win Window, err error) {
	// Initialize the video subsystem.
	if winCount == 0 {
		if C.SDL_InitSubSystem(C.SDL_INIT_VIDEO) != 0 {
			return Window{}, getLastError()
		}
	}

	// Open the window.
	var cFlags C.Uint32
	for _, flag := range flags {
		cFlags |= C.Uint32(flag)
	}
	title := C.CString("untitled")
	x := C.int(C.SDL_WINDOWPOS_UNDEFINED)
	y := C.int(C.SDL_WINDOWPOS_UNDEFINED)
	win.win = C.SDL_CreateWindow(title, x, y, C.int(width), C.int(height), cFlags)
	if win.win == nil {
		return Window{}, getLastError()
	}
	winCount++

	return win, nil
}

// Close closes the window.
func (win Window) Close() {
	C.SDL_DestroyWindow(win.win)
	winCount--

	// Terminate the video subsystem.
	if winCount == 0 {
		C.SDL_QuitSubSystem(C.SDL_INIT_VIDEO)
	}
}

// SetTitle sets the title of the window.
func (win Window) SetTitle(title string) {
	panic("not yet implemented")
}

// ShowCursor displays or hides the mouse cursor depending on the value of
// visible. It is visible by default.
func (win Window) ShowCursor(visible bool) {
	panic("not yet implemented")
}

// Width returns the width of the window.
func (win Window) Width() int {
	panic("not yet implemented")
}

// Height returns the height of the window.
func (win Window) Height() int {
	panic("not yet implemented")
}

// Draw draws the entire src image onto the window starting at the destination
// point dp.
func (win Window) Draw(dp image.Point, src wandi.Image) (err error) {
	panic("not yet implemented")
}

// DrawRect draws a subset of the src image, as defined by the source rectangle
// sr, onto the window starting at the destination point dp.
func (win Window) DrawRect(dp image.Point, src wandi.Image, sr image.Rectangle) (err error) {
	panic("not yet implemented")
}

// Fill fills the entire window with the provided color.
func (win Window) Fill(c color.Color) {
	panic("not yet implemented")
}

// Display displays what has been rendered so far to the window.
func (win Window) Display() {
	panic("not yet implemented")
}
