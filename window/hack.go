package window

// #include <SDL2/SDL.h>
import "C"

import (
	"unsafe"

	"github.com/mewmew/sdl/texture"
)

// drawableHack is a copy of texture.Drawable without modifications. Through the
// use of unsafe and with knowledge of its memory layout we are able to access
// unexported members. This hack allows us to cross package barriers while
// keeping the exported API clean.
type drawableHack struct {
	// A drawable GPU texture.
	tex *C.SDL_Texture
}

// drawableTexture returns the texture of the provided texture.Drawable.
func drawableTexture(tex texture.Drawable) *C.SDL_Texture {
	return (*drawableHack)(unsafe.Pointer(&tex)).tex
}

// imageHack is a copy of texture.Image without modifications. Through the use
// of unsafe and with knowledge of its memory layout we are able to access
// unexported members. This hack allows us to cross package barriers while
// keeping the exported API clean.
type imageHack struct {
	// A read-only GPU texture.
	tex *C.SDL_Texture
}

// imageTexture returns the texture of the provided texture.Image.
func imageTexture(tex texture.Image) *C.SDL_Texture {
	return (*imageHack)(unsafe.Pointer(&tex)).tex
}
