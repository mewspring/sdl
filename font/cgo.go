package font

// #cgo pkg-config: SDL2_ttf
// #include <SDL2/SDL_ttf.h>
import "C"

import (
	"errors"
	"image/color"
)

// cColor converts a Go color.Color to a C SDL_Color.
func cColor(c color.Color) (cColor C.SDL_Color) {
	r, g, b, a := c.RGBA()
	cColor.r = C.Uint8(r)
	cColor.g = C.Uint8(g)
	cColor.b = C.Uint8(b)
	cColor.a = C.Uint8(a)
	return cColor
}

// getTTFError returns the last TTF error message.
func getTTFError() (err error) {
	return errors.New(C.GoString(C.TTF_GetError()))
}

// getSDLError returns the last SDL error message.
func getSDLError() (err error) {
	return errors.New(C.GoString(C.SDL_GetError()))
}
