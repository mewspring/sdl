// Package render provides support for 2D accelerated rendering. At least one
// active rendering context is required for texture creation and drawing
// operations. A call to window.Open will implicitly add its renderer to the
// list of active renderers, and a subsequent call to win.Close will remove said
// renderer.
package render

// #cgo pkg-config: sdl2
// #include <SDL2/SDL.h>
import "C"

import (
	"errors"
	"unsafe"
)

// active keeps track of all active renderers.
var active = make(map[unsafe.Pointer]bool)

// Add adds the renderer to the list of active renderers.
func Add(renderer unsafe.Pointer) {
	if !active[renderer] {
		active[renderer] = true
	}
}

// Del deletes the renderer from the list of active renderers.
func Del(renderer unsafe.Pointer) {
	delete(active, renderer)
}

// getActive returns an active renderer or an error if there are no active
// renderers.
func getActive() (ren *C.SDL_Renderer, err error) {
	for renderer := range active {
		// TODO(u): Verify if switching between two different rendering contexts
		// too often affects performance.
		return (*C.SDL_Renderer)(renderer), nil
	}
	return nil, errors.New("render.getActive: unable to locate an active rendering context")
}

// createTexture returns a new texture of the specified dimensions, which may be
// the target of drawing operations if drawable is set to true.
func createTexture(width, height int, drawable bool) (tex unsafe.Pointer, err error) {
	ren, err := getActive()
	if err != nil {
		return nil, err
	}
	format := C.Uint32(C.SDL_PIXELFORMAT_ABGR8888)
	access := C.int(C.SDL_TEXTUREACCESS_STATIC)
	if drawable {
		access = C.SDL_TEXTUREACCESS_TARGET
	}
	tex = unsafe.Pointer(C.SDL_CreateTexture(ren, format, access, C.int(width), C.int(height)))
	if tex == nil {
		return nil, getLastError()
	}
	return tex, nil
}

// NewImage returns a new read-only texture of the specified dimensions.
func NewImage(width, height int) (tex unsafe.Pointer, err error) {
	return createTexture(width, height, false)
}

// NewDrawable returns a new drawable texture of the specified dimensions.
func NewDrawable(width, height int) (tex unsafe.Pointer, err error) {
	return createTexture(width, height, true)
}
