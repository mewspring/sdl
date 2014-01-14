// Package win handles window creation, drawing and events. The window events
// are defined in a dedicated package located at:
//    github.com/mewmew/we
//
// The library uses a small subset of the features provided by SDL version 2.0.
// For the sake of simplicity support for multiple windows has intentionally
// been left out.
package win

// #cgo pkg-config: sdl2
// #include <string.h>
// #include <SDL2/SDL.h>
//
// static SDL_Rect * makeRectArray(int n) {
//    return calloc(n, sizeof(SDL_Rect));
// }
//
// static void setArrayRect(SDL_Rect *rects, SDL_Rect *rect, int i) {
//    void *dst = rects + i*sizeof(SDL_Rect);
//    memcpy(dst, rect, sizeof(SDL_Rect));
// }
import "C"

import (
	"image"
	"unsafe"
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

// w represents the graphics window which is opened through a call to Open. It
// is this single window that is utilized throughout the entire library.
var w *C.SDL_Window

// Open opens a window with the specified dimensions and optional window flags.
// Only one window can be open at the same time. It is this single window that
// is utilized throughout the entire library. By default the window is not
// resizeable.
//
// Note: The Close function must be called when finished using the window.
func Open(width, height int, flags ...WindowFlag) (err error) {
	if w != nil {
		panic("win.Open: the window has already been opened.")
	}
	// Initialize the SDL video subsystem.
	if C.SDL_Init(C.SDL_INIT_VIDEO) != 0 {
		return getError()
	}

	// Open the window.
	var cFlags C.Uint32
	for _, flag := range flags {
		cFlags |= C.Uint32(flag)
	}
	title := C.CString("untitled")
	x := C.int(C.SDL_WINDOWPOS_UNDEFINED)
	y := C.int(C.SDL_WINDOWPOS_UNDEFINED)
	w = C.SDL_CreateWindow(title, x, y, C.int(width), C.int(height), cFlags)
	if w == nil {
		return getError()
	}
	return nil
}

// Close closes the window.
func Close() {
	C.SDL_DestroyWindow(w)
	w = nil
	C.SDL_Quit()
}

// SetTitle sets the title of the window.
func SetTitle(title string) {
	C.SDL_SetWindowTitle(w, C.CString(title))
}

// Screen returns the image associated with the window.
func Screen() (screen *Image, err error) {
	screen = new(Image)
	screen.s = C.SDL_GetWindowSurface(w)
	if screen.s == nil {
		return nil, getError()
	}
	screen.Width = int(screen.s.w)
	screen.Height = int(screen.s.h)
	return screen, nil
}

// Update copies the entire window image onto the screen.
func Update() (err error) {
	if C.SDL_UpdateWindowSurface(w) != 0 {
		return getError()
	}
	return nil
}

// UpdateRects copies a portion of the window image onto the screen as specified
// by rects.
func UpdateRects(rects []image.Rectangle) (err error) {
	cRects := C.makeRectArray(C.int(len(rects)))
	defer C.SDL_free(unsafe.Pointer(cRects))
	for i, rect := range rects {
		cRect := cRect(rect)
		C.setArrayRect(cRects, cRect, C.int(i))
	}
	if C.SDL_UpdateWindowSurfaceRects(w, cRects, C.int(len(rects))) != 0 {
		return getError()
	}
	return nil
}

// Draw draws the entire src image onto the window starting at the destination
// point dp.
func Draw(dp image.Point, src *Image) (err error) {
	dst, err := Screen()
	if err != nil {
		return err
	}
	return dst.Draw(dp, src)
}

// DrawRect fills the destination rectangle dr of the window with corresponding
// pixels from the src image starting at the source point sp.
func DrawRect(dr image.Rectangle, src *Image, sp image.Point) (err error) {
	dst, err := Screen()
	if err != nil {
		return err
	}
	return dst.DrawRect(dr, src, sp)
}
