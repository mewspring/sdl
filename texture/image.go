package texture

// #cgo pkg-config: sdl2
// #include <SDL2/SDL.h>
import "C"

import (
	"image"
	"log"
	"unsafe"

	"github.com/mewmew/sdl"
)

// Image represent a read-only texture. It implements the wandi.Image interface.
type Image struct {
	// A read-only GPU texture.
	tex *C.SDL_Texture
}

// Load loads the provided file and converts it into a read-only texture.
//
// Note: The Free method of the texture must be called when finished using it.
func Load(path string) (tex Image, err error) {
	panic("not yet implemented")
}

// Read reads the provided image and converts it into a read-only texture.
//
// Note: The Free method of the texture must be called when finished using it.
func Read(src image.Image) (tex Image, err error) {
	// Create a texture of the same dimensions as the provided image.
	bounds := src.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	format := C.Uint32(C.SDL_PIXELFORMAT_ABGR8888)
	access := C.int(C.SDL_TEXTUREACCESS_STATIC)
	tex.tex = C.SDL_CreateTexture(ren, format, access, C.int(width), C.int(height))
	if tex.tex == nil {
		return Image{}, getLastError()
	}

	// Copy the image's pixels to the texture.
	rgba, ok := src.(*image.RGBA)
	if !ok {
		log.Fatalf("image format %T not yet supported.\n", src)
	}
	pix := unsafe.Pointer(&rgba.Pix[0])
	if C.SDL_UpdateTexture(tex.tex, nil, pix, C.int(rgba.Stride)) != 0 {
		return Image{}, getLastError()
	}

	return tex, nil
}

// Free frees the texture.
func (tex Image) Free() {
	C.SDL_DestroyTexture(tex.tex)
}

// Width returns the width of the texture.
func (tex Image) Width() int {
	panic("not yet implemented")
}

// Height returns the height of the texture.
func (tex Image) Height() int {
	panic("not yet implemented")
}
