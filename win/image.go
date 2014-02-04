package win

// #cgo pkg-config: sdl2
// #include <string.h>
// #include <SDL2/SDL.h>
import "C"

import (
	"image"
	"log"
	"unsafe"

	"github.com/mewkiz/pkg/imgutil"
)

// An Image is a mutable collection of pixels.
type Image struct {
	// The width and height of the image.
	Width, Height int
	// C texture pointer.
	tex *C.SDL_Texture
}

// NewImage returns a new image of the specified dimensions.
//
// Note: The Free method of the image should be called when finished using it.
func NewImage(width, height int) (img *Image, err error) {
	access := C.int(C.SDL_TEXTUREACCESS_TARGET)
	tex := C.SDL_CreateTexture(renderer, C.SDL_PIXELFORMAT_ABGR8888, access, C.int(width), C.int(height))
	if tex == nil {
		return nil, getSDLError()
	}
	if C.SDL_SetTextureBlendMode(tex, C.SDL_BLENDMODE_BLEND) != 0 {
		return nil, getSDLError()
	}
	img = &Image{
		Width:  width,
		Height: height,
		tex:    tex,
	}
	return img, nil
}

// LoadImage loads the provided image file and returns it as an image.
//
// Note: The Free method of the image should be called when finished using it.
func LoadImage(imgPath string) (img *Image, err error) {
	src, err := imgutil.ReadFile(imgPath)
	if err != nil {
		return nil, err
	}
	return ReadImage(src)
}

// ReadImage reads the provided image, converts it to the standard image format
// of this library and returns it.
//
// Note: The Free method of the image should be called when finished using it.
func ReadImage(src image.Image) (img *Image, err error) {
	bounds := src.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	img, err = NewImage(width, height)
	if err != nil {
		return nil, err
	}
	dstRect := cRect(image.Rect(0, 0, width, height))
	switch i := src.(type) {
	case *image.NRGBA:
		C.SDL_UpdateTexture(img.tex, dstRect, unsafe.Pointer(&i.Pix[0]), C.int(i.Stride))
	case *image.RGBA:
		// TODO(u): Do we need normalize the image since it is stored as
		// premultiplied alpha?
		C.SDL_UpdateTexture(img.tex, dstRect, unsafe.Pointer(&i.Pix[0]), C.int(i.Stride))
	default:
		log.Fatalf("win.ReadImage: image format %T not yet supported.\n", i)
	}
	return img, nil
}

// Free frees the image.
func (img *Image) Free() {
	C.SDL_DestroyTexture(img.tex)
}

// Draw draws the entire src image onto the dst image starting at the
// destination point dp.
func (dst *Image) Draw(dp image.Point, src *Image) (err error) {
	dr := image.Rect(dp.X, dp.Y, dp.X+src.Width, dp.Y+src.Height)
	return dst.DrawRect(dr, src, image.ZP)
}

// DrawRect fills the destination rectangle dr of the dst image with
// corresponding pixels from the src image starting at the source point sp.
func (dst *Image) DrawRect(dr image.Rectangle, src *Image, sp image.Point) (err error) {
	// Set the texture as the current rendering target.
	if C.SDL_SetRenderTarget(renderer, dst.tex) != 0 {
		return getSDLError()
	}

	sr := image.Rect(sp.X, sp.Y, sp.X+dr.Dx(), sp.Y+dr.Dy())
	srcRect := cRect(sr)
	dstRect := cRect(dr)
	if C.SDL_RenderCopy(renderer, src.tex, srcRect, dstRect) != 0 {
		return getSDLError()
	}
	return nil
}
