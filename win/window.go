// Package win handles window creation, event handling and image drawing.
//
// The library uses a small subset of the features provided by the SDL library
// version 2.0 [1].
//
// [1]: http://libsdl.org/
package win

// #cgo pkg-config: sdl2
// #include <SDL2/SDL.h>
import "C"

import (
	"image"
	"image/color"

	"github.com/mewmew/wandi"
)

// WindowFlag is a bitfield of window flags.
type WindowFlag uint32

// Window flags.
const (
	// Resizeable states that the window can be resized.
	Resizeable WindowFlag = C.SDL_WINDOW_RESIZABLE
	// FullScreen states that the window is in full screen mode.
	FullScreen WindowFlag = C.SDL_WINDOW_FULLSCREEN
)

// A sdlWindow represents a graphical window capable of handling draw operations
// and window events.
type sdlWindow struct {
	// Pointer to the C SDL_Window of the window.
	w *C.SDL_Window
}

// Open opens a window with the specified dimensions and optional window flags.
// By default the window is not resizeable.
//
// Note: The Close method of the window must be called when finished using it.
func Open(width, height int, flags ...WindowFlag) (win wandi.Window, err error) {
	// Initialize the SDL video subsystem.
	if C.SDL_Init(C.SDL_INIT_VIDEO) != 0 {
		return nil, getSDLError()
	}

	// Open the window.
	var cFlags C.Uint32
	for _, flag := range flags {
		cFlags |= C.Uint32(flag)
	}
	title := C.CString("untitled")
	x := C.int(C.SDL_WINDOWPOS_UNDEFINED)
	y := C.int(C.SDL_WINDOWPOS_UNDEFINED)
	sdlWin := new(sdlWindow)
	sdlWin.w = C.SDL_CreateWindow(title, x, y, C.int(width), C.int(height), cFlags)
	if sdlWin.w == nil {
		return nil, getSDLError()
	}

	// Make sure the window surface is valid for updates.
	s := C.SDL_GetWindowSurface(sdlWin.w)
	if s == nil {
		return nil, getSDLError()
	}

	return sdlWin, nil
}

// Close closes the window.
func (sdlWin *sdlWindow) Close() {
	C.SDL_DestroyWindow(sdlWin.w)
	C.SDL_Quit()
}

// SetTitle sets the title of the window.
func (sdlWin *sdlWindow) SetTitle(title string) {
	C.SDL_SetWindowTitle(sdlWin.w, C.CString(title))
}

// Screen returns the image associated with the window.
func (sdlWin *sdlWindow) Screen() (screen wandi.Image, err error) {
	return sdlWin.screen()
}

// screen returns the image associated with the window.
func (sdlWin *sdlWindow) screen() (sdlScreen *Image, err error) {
	sdlScreen = new(Image)
	sdlScreen.s = C.SDL_GetWindowSurface(sdlWin.w)
	if sdlScreen.s == nil {
		return nil, getSDLError()
	}
	sdlScreen.w = int(sdlScreen.s.w)
	sdlScreen.h = int(sdlScreen.s.h)
	return sdlScreen, nil
}

// Update copies the entire window image onto the screen.
func (sdlWin *sdlWindow) Update() (err error) {
	if C.SDL_UpdateWindowSurface(sdlWin.w) != 0 {
		return getSDLError()
	}
	return nil
}

func (sdlWin *sdlWindow) Clear(c color.Color) (err error) {
	panic("sdlWindow.Clear: not yet implemented.")
}

// Draw draws the entire src image onto the window starting at the destination
// point dp.
func (sdlWin *sdlWindow) Draw(dp image.Point, src wandi.Image) (err error) {
	dst, err := sdlWin.screen()
	if err != nil {
		return err
	}
	return dst.Draw(dp, src)
}

// DrawRect fills the destination rectangle dr of the window with corresponding
// pixels from the src image starting at the source point sp.
func (sdlWin *sdlWindow) DrawRect(dr image.Rectangle, src wandi.Image, sp image.Point) (err error) {
	dst, err := sdlWin.screen()
	if err != nil {
		return err
	}
	return dst.DrawRect(dr, src, sp)
}
