package win

// #cgo pkg-config: sdl2
// #include <SDL2/SDL.h>
import "C"

import (
	"errors"
	"image"
)

// cRect converts a Go image.Rectangle to a C SDL_Rect.
func cRect(rect image.Rectangle) (cRect *C.SDL_Rect) {
	if rect == image.ZR {
		return nil
	}
	cRect = new(C.SDL_Rect)
	cRect.x = C.int(rect.Min.X)
	cRect.y = C.int(rect.Min.Y)
	cRect.w = C.int(rect.Max.X - rect.Min.X)
	cRect.h = C.int(rect.Max.Y - rect.Min.Y)
	return cRect
}

// getSDLError returns the last SDL error message.
func getSDLError() (err error) {
	return errors.New(C.GoString(C.SDL_GetError()))
}
