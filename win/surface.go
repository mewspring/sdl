package win

// #cgo pkg-config: sdl2
// #include <SDL2/SDL.h>
import "C"

import (
	"encoding/binary"
	"image"
	"image/draw"
	"reflect"
	"unsafe"

	"github.com/mewkiz/pkg/imgutil"
)

// An Image is a collection of pixels.
type Image struct {
	// Width and height of the image.
	Width, Height int
	// C surface pointer.
	s *C.SDL_Surface
}

// NewImage returns a new image of the specified dimensions.
//
// Note: The Free method of the image should be called when finished using it.
func NewImage(width, height int) (img *Image, err error) {
	img = &Image{
		Width:  width,
		Height: height,
	}
	// Red, green blue and alpha masks.
	var r, g, b, a C.Uint32
	switch nativeByteOrder {
	case binary.LittleEndian:
		r = 0x000000FF
		g = 0x0000FF00
		b = 0x00FF0000
		a = 0xFF000000
	case binary.BigEndian:
		r = 0xFF000000
		g = 0x00FF0000
		b = 0x0000FF00
		a = 0x000000FF
	}
	img.s = C.SDL_CreateRGBSurface(0, C.int(width), C.int(height), 32, r, g, b, a)
	if img.s == nil {
		return nil, getError()
	}
	return img, nil
}

// LoadImage loads and returns the provided image.
//
// Note: The Free method of the image should be called when finished using it.
func LoadImage(imgPath string) (img *Image, err error) {
	src, err := imgutil.ReadFile(imgPath)
	if err != nil {
		return nil, err
	}
	return ConvertImage(src)
}

// ConvertImage converts the provided image to the standard image format of this
// library.
//
// Note: The Free method of the image should be called when finished using it.
func ConvertImage(src image.Image) (img *Image, err error) {
	rect := src.Bounds()
	width, height := rect.Dx(), rect.Dy()
	img, err = NewImage(width, height)
	if err != nil {
		return nil, err
	}
	var pix []uint8
	switch i := src.(type) {
	case *image.NRGBA:
		pix = i.Pix
	case *image.RGBA:
		// TODO(u): Do we need normalize the image since its stored as
		// premultiplied alpha. If so divide the color values by the alpha value.
		// An alternative is to use the default fallback. If so benchmark first.
		pix = i.Pix
	default:
		copyPixels(img, src)
		return img, nil
	}
	C.memcpy(img.s.pixels, unsafe.Pointer(&pix[0]), C.size_t(len(pix)))
	return img, nil
}

// copyPixels provides a draw.Draw fallback for arbitrary image formats. It
// draws directly to the memory of the SDL surface using unsafe.
//
// Note: The dst image must be a valid SDL surface created with NewImage.
func copyPixels(dst *Image, src image.Image) {
	// stride is the size in bytes of each line. The size of an individual
	// pixel is 4 bytes.
	stride := dst.Width * 4
	// size is the total size in bytes of the pixel data.
	size := dst.Height * stride
	sh := reflect.SliceHeader{
		Data: uintptr(dst.s.pixels),
		Len:  size,
		Cap:  size,
	}
	// dstPix is a byte slice which points to the memory of the surface's pixels.
	dstPix := *(*[]byte)(unsafe.Pointer(&sh))

	dstRect := image.Rect(0, 0, dst.Width, dst.Height)
	dstImg := &image.NRGBA{
		Pix:    dstPix,
		Stride: stride,
		Rect:   dstRect,
	}
	draw.Draw(dstImg, dstRect, src, image.ZP, draw.Over)
}

// Free frees the image.
func (img *Image) Free() {
	C.SDL_FreeSurface(img.s)
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
	sr := image.Rect(sp.X, sp.Y, sp.X+dr.Dx(), sp.Y+dr.Dy())
	srcRect := cRect(sr)
	dstRect := cRect(dr)
	if C.SDL_BlitSurface(src.s, srcRect, dst.s, dstRect) != 0 {
		return getError()
	}
	return nil
}
