package window

// #cgo pkg-config: sdl2
// #include <string.h>
// #include <SDL2/SDL.h>
import "C"

import (
	"fmt"
	"image"
	"image/draw"
	"reflect"
	"unsafe"

	"github.com/mewkiz/pkg/imgutil"
	"github.com/mewmew/wandi"
)

// An Image is a mutable collection of pixels, whose memory can be freed. It
// implements the wandi.Image interface.
type Image struct {
	// The width and height of the image.
	w, h int
	// Pointer to the C SDL_Surface of the image.
	s *C.SDL_Surface
}

// NewImage returns a new image of the specified dimensions.
//
// Note: The Free method of the image should be called when finished using it.
func NewImage(width, height int) (img *Image, err error) {
	return newImage(width, height)
}

// newImage returns a new image of the specified dimensions.
//
// Note: The Free method of the image should be called when finished using it.
func newImage(width, height int) (sdlImg *Image, err error) {
	sdlImg = &Image{
		w: width,
		h: height,
	}
	// Red, green blue and alpha masks.
	var r, g, b, a C.Uint32
	if nativeBigEndian {
		r = 0xFF000000
		g = 0x00FF0000
		b = 0x0000FF00
		a = 0x000000FF
	} else {
		r = 0x000000FF
		g = 0x0000FF00
		b = 0x00FF0000
		a = 0xFF000000
	}
	sdlImg.s = C.SDL_CreateRGBSurface(0, C.int(width), C.int(height), 32, r, g, b, a)
	if sdlImg.s == nil {
		return nil, getSDLError()
	}
	return sdlImg, nil
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
	rect := src.Bounds()
	width, height := rect.Dx(), rect.Dy()
	sdlImg, err := newImage(width, height)
	if err != nil {
		return nil, err
	}

	var pix []uint8
	switch i := src.(type) {
	case *image.NRGBA:
		pix = i.Pix
	case *image.RGBA:
		// TODO(u): Do we need normalize the image since its stored as
		// premultiplied alpha? If so divide the color values by the alpha value.
		// An alternative is to use the default fallback. If so benchmark first.
		pix = i.Pix
	default:
		copyPixels(sdlImg, src)
		return sdlImg, nil
	}
	C.memcpy(sdlImg.s.pixels, unsafe.Pointer(&pix[0]), C.size_t(len(pix)))
	return sdlImg, nil
}

// copyPixels copies the pixels of the src image to the dst SDL surface. It uses
// unsafe to draw directly to the memory of the SDL surface. No alpha blending
// is performed since it's always used during the creation of new SDL surfaces.
//
// Note: The dst image must be a valid SDL surface created with NewImage.
func copyPixels(dst *Image, src image.Image) {
	// stride is the size in bytes of each line. The size of an individual
	// pixel is 4 bytes.
	stride := dst.w * 4
	// size is the total size in bytes of the pixel data.
	size := dst.h * stride
	sh := reflect.SliceHeader{
		Data: uintptr(dst.s.pixels),
		Len:  size,
		Cap:  size,
	}
	// dstPix is a byte slice which points to the memory of the surface's pixels.
	dstPix := *(*[]byte)(unsafe.Pointer(&sh))

	dstRect := image.Rect(0, 0, dst.w, dst.h)
	dstImg := &image.NRGBA{
		Pix:    dstPix,
		Stride: stride,
		Rect:   dstRect,
	}
	draw.Draw(dstImg, dstRect, src, image.ZP, draw.Src)
}

// Free frees the image.
func (sdlImg *Image) Free() {
	C.SDL_FreeSurface(sdlImg.s)
}

// Width returns the width of the image.
func (sdlImg *Image) Width() int {
	return sdlImg.w
}

// Height returns the height of the image.
func (sdlImg *Image) Height() int {
	return sdlImg.h
}

// Draw draws the entire src image onto the dst image starting at the
// destination point dp.
func (dst *Image) Draw(dp image.Point, src wandi.Image) (err error) {
	sr := image.Rect(0, 0, src.Width(), src.Height())
	return dst.DrawRect(dp, src, sr)
}

// DrawRect draws a subset of the src image, as defined by the source rectangle
// sr, onto the dst image starting at the destination point dp.
func (dst *Image) DrawRect(dp image.Point, src wandi.Image, sr image.Rectangle) (err error) {
	switch srcImg := src.(type) {
	case *Image:
		dr := image.Rect(dp.X, dp.Y, dp.X+sr.Dx(), dp.Y+sr.Dy())
		srcRect := cRect(sr)
		dstRect := cRect(dr)
		if C.SDL_BlitSurface(srcImg.s, srcRect, dst.s, dstRect) != 0 {
			return getSDLError()
		}
		return nil
	default:
		return fmt.Errorf("Image.DrawRect: unsupported image type %T", src)
	}
}
