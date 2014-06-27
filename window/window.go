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

// active represent the number of active windows.
var active int

// Open opens a new window of the specified dimensions and optional window
// flags. By default the window is not resizeable.
//
// Note: The Close method of the window must be called when finished using it.
func Open(width, height int, flags ...Flag) (win Window, err error) {
	// Initialize the video subsystem.
	if active == 0 {
		if C.SDL_InitSubSystem(C.SDL_INIT_VIDEO) != 0 {
			return Window{}, fmt.Errorf("window.Open: %v", getLastError())
		}
	}
	active++

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
		return Window{}, fmt.Errorf("window.Open: %v", getLastError())
	}

	// Create a renderer for the window.
	win.ren = C.SDL_CreateRenderer(win.win, -1, C.SDL_RENDERER_PRESENTVSYNC)
	if win.ren == nil {
		return Window{}, fmt.Errorf("window.Open: %v", getLastError())
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
	active--
	if active == 0 {
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
func (win Window) ShowCursor(visible bool) (err error) {
	toggle := C.int(0)
	if visible {
		toggle = 1
	}
	if C.SDL_ShowCursor(toggle) < 0 {
		return fmt.Errorf("Window.ShowCursor: %v", getLastError())
	}
	return nil
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

// drawRect draws a subset of the src texture, as defined by the source
// rectangle sr, onto the window starting at the destination point dp.
func drawRect(ren *C.SDL_Renderer, dp image.Point, src *C.SDL_Texture, sr image.Rectangle) (err error) {
	width, height := C.int(sr.Dx()), C.int(sr.Dy())
	srcrect := &C.SDL_Rect{
		x: C.int(sr.Min.X),
		y: C.int(sr.Min.Y),
		w: width,
		h: height,
	}
	dstrect := &C.SDL_Rect{
		x: C.int(dp.X),
		y: C.int(dp.Y),
		w: width,
		h: height,
	}
	if C.SDL_RenderCopy(ren, src, srcrect, dstrect) != 0 {
		return fmt.Errorf("window.drawRect: %v", getLastError())
	}
	return nil
}

// DrawRect draws a subset of the src image, as defined by the source rectangle
// sr, onto the window starting at the destination point dp.
func (win Window) DrawRect(dp image.Point, src wandi.Image, sr image.Rectangle) (err error) {
	switch srcImg := src.(type) {
	case texture.Drawable:
		srcTex := drawableTexture(srcImg)
		return drawRect(win.ren, dp, srcTex, sr)
	case texture.Image:
		srcTex := imageTexture(srcImg)
		return drawRect(win.ren, dp, srcTex, sr)
	default:
		return fmt.Errorf("Window.DrawRect: source type %T not yet supported", src)
	}
}

// Fill fills the entire window with the provided color.
func (win Window) Fill(c color.Color) (err error) {
	r, g, b, a := c.RGBA()
	if C.SDL_SetRenderDrawColor(win.ren, C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)) != 0 {
		return fmt.Errorf("Window.Fill: %v", getLastError())
	}
	if C.SDL_RenderClear(win.ren) != 0 {
		return fmt.Errorf("Window.Fill: %v", getLastError())
	}
	return nil
}

// Display displays what has been rendered so far to the window.
func (win Window) Display() {
	C.SDL_RenderPresent(win.ren)
}
