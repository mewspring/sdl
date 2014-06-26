package texture

// #cgo pkg-config: sdl2
// #include <SDL2/SDL.h>
import "C"

import (
	"fmt"
	"image"
	"image/draw"
	"log"
	"time"
	"unsafe"

	"github.com/mewkiz/pkg/imgutil"
)

// Image represent a read-only texture. It implements the wandi.Image interface.
type Image struct {
	// A read-only GPU texture.
	tex *C.SDL_Texture
}

// newImage creates a read-only texture of the specified dimensions.
func newImage(width, height int) (tex Image, err error) {
	tex.tex, err = create(width, height, false)
	if err != nil {
		return Image{}, err
	}
	return tex, nil
}

// Load loads the provided file and converts it into a read-only texture.
//
// Note: The Free method of the texture must be called when finished using it.
func Load(path string) (tex Image, err error) {
	src, err := imgutil.ReadFile(path)
	if err != nil {
		return Image{}, err
	}
	return Read(src)
}

// fallback converts the provided image or subimage into a RGBA image.
func fallback(src image.Image) *image.RGBA {
	start := time.Now()

	// Create a new RGBA image and draw the src image onto it.
	bounds := src.Bounds()
	dr := image.Rect(0, 0, bounds.Dx(), bounds.Dy())
	dst := image.NewRGBA(dr)
	draw.Draw(dst, dr, src, bounds.Min, draw.Src)

	log.Printf("texture.fallback: fallback conversion for non-RGBA image (%T) finished in: %v.\n", src, time.Since(start))

	return dst
}

// Read reads the provided image and converts it into a read-only texture.
//
// Note: The Free method of the texture must be called when finished using it.
func Read(src image.Image) (tex Image, err error) {
	// Use fallback conversion for unknown image formats.
	rgba, ok := src.(*image.RGBA)
	if !ok {
		return Read(fallback(src))
	}

	// Use fallback conversion for subimages.
	width, height := rgba.Rect.Dx(), rgba.Rect.Dy()
	if rgba.Stride != 4*width {
		return Read(fallback(src))
	}

	// Create a read-only texture based on the pixels of the src image.
	tex, err = newImage(width, height)
	if err != nil {
		return Image{}, err
	}
	pix := unsafe.Pointer(&rgba.Pix[0])
	if C.SDL_UpdateTexture(tex.tex, nil, pix, C.int(rgba.Stride)) != 0 {
		tex.Free()
		return Image{}, fmt.Errorf("texture.Read: %v", getLastError())
	}

	return tex, nil
}

// Free frees the texture.
func (tex Image) Free() {
	C.SDL_DestroyTexture(tex.tex)
}

// Width returns the width of the texture.
func (tex Image) Width() int {
	var width C.int
	if C.SDL_QueryTexture(tex.tex, nil, nil, &width, nil) != 0 {
		log.Fatalln("Image.Width: unable to locate texture width;", getLastError())
	}
	return int(width)
}

// Height returns the height of the texture.
func (tex Image) Height() int {
	var height C.int
	if C.SDL_QueryTexture(tex.tex, nil, nil, nil, &height) != 0 {
		log.Fatalln("Image.Height: unable to locate texture height;", getLastError())
	}
	return int(height)
}
