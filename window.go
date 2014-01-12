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

// A Window represents a single graphics window.
type Window struct {
	// C window pointer.
	w *C.SDL_Window
}

// TODO(u): make it possible to create resizeable and fullscreen windows.

// OpenWindow opens a new window of the specified dimensions.
//
// Note: The Close method of the window should be called when finished using it.
func OpenWindow(width, height int) (w *Window, err error) {
	w = new(Window)
	title := C.CString("untitled")
	x := C.int(C.SDL_WINDOWPOS_UNDEFINED)
	y := C.int(C.SDL_WINDOWPOS_UNDEFINED)
	w.w = C.SDL_CreateWindow(title, x, y, C.int(width), C.int(height), 0)
	if w.w == nil {
		return nil, getError()
	}
	return w, nil
}

// Close closes the window.
func (w *Window) Close() {
	C.SDL_DestroyWindow(w.w)
}

// SetTitle sets the title of the window.
func (w *Window) SetTitle(title string) {
	C.SDL_SetWindowTitle(w.w, C.CString(title))
}

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
