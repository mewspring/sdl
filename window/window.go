// Package window handles window creation, drawing and events. It uses a small
// subset of the features provided by the SDL library version 2.0 [1].
//
// [1]: http://www.libsdl.org/
package window

// #cgo pkg-config: sdl2
// #include <SDL2/SDL.h>
import "C"

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"unsafe"

	"github.com/mewmew/sdl/texture"
	"github.com/mewmew/wandi"
)

// Flag is a bitfield of window flags.
type Flag uint32

// Window flags.
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
	// The rendering context associated with the window.
	ren *C.SDL_Renderer
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
		// TODO(u): Add a goroutine which does event polling and sends the events
		// to their corresponding window.
	}
	winCount++

	// Open the window.
	var cFlags C.Uint32
	for _, flag := range flags {
		cFlags |= C.Uint32(flag)
	}
	title := C.CString("untitled")
	x := C.int(C.SDL_WINDOWPOS_UNDEFINED)
	y := C.int(C.SDL_WINDOWPOS_UNDEFINED)
	// TODO(u): Make use of SDL_SetHint and SDL_RenderSetLogicalSize for
	// full screen windows.
	win.win = C.SDL_CreateWindow(title, x, y, C.int(width), C.int(height), cFlags)
	if win.win == nil {
		return Window{}, getLastError()
	}

	// Create a renderer for the window.
	win.ren = C.SDL_CreateRenderer(win.win, -1, C.SDL_RENDERER_PRESENTVSYNC)
	if win.ren == nil {
		return Window{}, getLastError()
	}
	texture.AddRenderer(unsafe.Pointer(win.ren))

	return win, nil
}

// Close closes the window.
func (win Window) Close() {
	if win.ren != nil {
		texture.DelRenderer(unsafe.Pointer(win.ren))
		C.SDL_DestroyRenderer(win.ren)
	}
	if win.win != nil {
		C.SDL_DestroyWindow(win.win)
	}

	// Terminate the video subsystem.
	winCount--
	if winCount == 0 {
		C.SDL_QuitSubSystem(C.SDL_INIT_VIDEO)
	}
}

// SetTitle sets the title of the window.
func (win Window) SetTitle(title string) {
	C.SDL_SetWindowTitle(win.win, C.CString(title))
}

// Show displays or hides the window depending on the value of visible. It is
// visible by default.
func (win Window) Show(visible bool) {
	if visible {
		C.SDL_ShowWindow(win.win)
	} else {
		C.SDL_HideWindow(win.win)
	}
}

// ShowCursor displays or hides the mouse cursor depending on the value of
// visible. It is visible by default.
func (win Window) ShowCursor(visible bool) {
	toggle := C.int(0)
	if visible {
		toggle = 1
	}
	if C.SDL_ShowCursor(toggle) < 0 {
		log.Println(getLastError())
	}
}

// Width returns the width of the window.
func (win Window) Width() int {
	var width C.int
	C.SDL_GetWindowSize(win.win, &width, nil)
	return int(width)
}

// Height returns the height of the window.
func (win Window) Height() int {
	var height C.int
	C.SDL_GetWindowSize(win.win, nil, &height)
	return int(height)
}

// Draw draws the entire src image onto the window starting at the destination
// point dp.
func (win Window) Draw(dp image.Point, src wandi.Image) (err error) {
	sr := image.Rect(0, 0, win.Width(), win.Height())
	return win.DrawRect(dp, src, sr)
}

// DrawRect draws a subset of the src image, as defined by the source rectangle
// sr, onto the window starting at the destination point dp.
func (win Window) DrawRect(dp image.Point, src wandi.Image, sr image.Rectangle) (err error) {
	switch srcImg := src.(type) {
	case texture.Drawable:
		tex := drawableTexture(srcImg)
		C.SDL_RenderCopy(win.ren, tex, nil, nil)
	case texture.Image:
		tex := imageTexture(srcImg)
		C.SDL_RenderCopy(win.ren, tex, nil, nil)
	default:
		return fmt.Errorf("Window.DrawRect: source type %T not yet supported", src)
	}
	return nil
}

// Fill fills the entire window with the provided color.
func (win Window) Fill(c color.Color) {
	r, g, b, a := c.RGBA()
	if C.SDL_SetRenderDrawColor(win.ren, C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)) != 0 {
		log.Fatalf("Window.Fill: %v\n", getLastError())
	}
	if C.SDL_RenderClear(win.ren) != 0 {
		log.Fatalf("Window.Fill: %v\n", getLastError())
	}
}

// Display displays what has been rendered so far to the window.
func (win Window) Display() {
	C.SDL_RenderPresent(win.ren)
}
