package texture

// #cgo pkg-config: sdl2
// #include <SDL2/SDL.h>
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

// TODO(u): Protect active from race conditions.

// active keeps track of all active renderers.
var active = make(map[unsafe.Pointer]bool)

// AddRenderer adds the renderer to the list of active renderers. At least one
// active rendering context is required for texture creation and drawing
// operations. A call to window.Open will implicitly add its renderer to the
// list of active renderers, and a subsequent call to win.Close will remove said
// renderer.
func AddRenderer(renderer unsafe.Pointer) {
	if !active[renderer] {
		active[renderer] = true
	}
}

// DelRenderer deletes the renderer from the list of active renderers.
func DelRenderer(renderer unsafe.Pointer) {
	delete(active, renderer)
}

// getRenderer returns an active renderer or an error if there are no active
// renderers.
func getRenderer() (ren *C.SDL_Renderer, err error) {
	for renderer := range active {
		// TODO(u): Verify if switching between two different rendering contexts
		// too often affects performance.
		return (*C.SDL_Renderer)(renderer), nil
	}
	return nil, errors.New("texture.getRenderer: unable to locate an active rendering context")
}

// create returns a new texture of the specified dimensions, which may be the
// target of drawing operations if drawable is set to true.
func create(width, height int, drawable bool) (tex *C.SDL_Texture, err error) {
	ren, err := getRenderer()
	if err != nil {
		return nil, err
	}
	format := C.Uint32(C.SDL_PIXELFORMAT_ABGR8888)
	access := C.int(C.SDL_TEXTUREACCESS_STATIC)
	if drawable {
		access = C.SDL_TEXTUREACCESS_TARGET
	}
	tex = C.SDL_CreateTexture(ren, format, access, C.int(width), C.int(height))
	if tex == nil {
		return nil, fmt.Errorf("texture.create: %v", getLastError())
	}
	return tex, nil
}
