package font

// #cgo pkg-config: SDL2_ttf
// #include <SDL2/SDL_ttf.h>
import "C"

import (
	"errors"
	"image/color"
	"unsafe"

	"github.com/mewmew/sdl/window"
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

// getTTFError returns the last error message.
func getTTFError() (err error) {
	return errors.New(C.GoString(C.TTF_GetError()))
}

// image is identical to the window.Image structure. It is used in conjunction
// with unsafe to avoid having to export the SDL surface member of the
// window.Image structure.
type image struct {
	// The width and height of the image.
	w, h int
	// Pointer to the C SDL_Surface of the image.
	s *C.SDL_Surface
}

// winImage returns a *window.Image corresponding to the provided SDL surface.
func winImage(s *C.SDL_Surface) *window.Image {
	// The SDL surface of window.Image isn't exported and it shouldn't be.
	// Therefore we use unsafe to keep the API clean.
	img := &image{
		w: int(s.w),
		h: int(s.h),
		s: s,
	}
	return (*window.Image)(unsafe.Pointer(img))
}
