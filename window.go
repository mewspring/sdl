package sdl

// #cgo pkg-config: sdl2
// #include <SDL2/SDL.h>
//
// static SDL_Rect * makeRectArray(int size) {
//    return calloc(sizeof(SDL_Rect), size);
// }
//
// static void setArrayRect(SDL_Rect *cRects, SDL_Rect *rect, int i) {
//    void *dst = cRects + sizeof(SDL_Rect)*i;
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
	// The Resizeable flag states that the window can be resized.
	Resizeable WindowFlag = C.SDL_WINDOW_RESIZABLE
	// The FullScreen flag states that the window is in full screen mode.
	FullScreen WindowFlag = C.SDL_WINDOW_FULLSCREEN
)

// A Window represents a single graphics window.
type Window struct {
	// C window pointer.
	w *C.SDL_Window
}

// OpenWindow opens a new window of the specified dimensions and optional window
// flags. By default the window is not resizeable.
//
// Note: The Close method of the window should be called when finished using it.
func OpenWindow(width, height int, flags ...WindowFlag) (w *Window, err error) {
	// Initialize SDL video subsystem.
	if C.SDL_Init(C.SDL_INIT_VIDEO) != 0 {
		return nil, getError()
	}

	// Open a new window.
	var cFlags C.Uint32
	for _, flag := range flags {
		cFlags |= C.Uint32(flag)
	}
	w = new(Window)
	title := C.CString("untitled")
	x := C.int(C.SDL_WINDOWPOS_UNDEFINED)
	y := C.int(C.SDL_WINDOWPOS_UNDEFINED)
	w.w = C.SDL_CreateWindow(title, x, y, C.int(width), C.int(height), cFlags)
	if w.w == nil {
		return nil, getError()
	}

	return w, nil
}

// Close closes the window.
func (w *Window) Close() {
	C.SDL_DestroyWindow(w.w)
	// TODO(u): Only quit the video subsystem if audio support is ever
	// implemented.
	C.SDL_Quit()
}

// SetTitle sets the title of the window.
func (w *Window) SetTitle(title string) {
	C.SDL_SetWindowTitle(w.w, C.CString(title))
}

// TODO(u): The sdl package has no intention of providing support for multiple
// windows. Think about how this could affect the API. For instance OpenWindow
// could initialize an unexported global window w. This would allow dedicated
// window drawing functions to be implemented which call w.Surface() to get
// access to the dst surface. A Surface method could still be usefull to provide
// screenshot functionality.

// Surface returns the surface associated with the window.
func (w *Window) Surface() (s *Surface, err error) {
	s = new(Surface)
	s.s = C.SDL_GetWindowSurface(w.w)
	if s.s == nil {
		return nil, getError()
	}
	return s, nil
}

// Update copies the entire window surface onto screen.
func (w *Window) Update() (err error) {
	if C.SDL_UpdateWindowSurface(w.w) != 0 {
		return getError()
	}
	return nil
}

// UpdateRects copies a portion of the window surface onto screen as specified
// by rects.
func (w *Window) UpdateRects(rects []image.Rectangle) (err error) {
	cRects := C.makeRectArray(C.int(len(rects)))
	defer C.SDL_free(unsafe.Pointer(cRects))
	for i, rect := range rects {
		cRect := cRect(rect)
		C.setArrayRect(cRects, cRect, C.int(i))
	}
	if C.SDL_UpdateWindowSurfaceRects(w.w, cRects, C.int(len(rects))) != 0 {
		return getError()
	}
	return nil
}
