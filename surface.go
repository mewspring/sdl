package sdl

// #cgo pkg-config: sdl2
// #include <SDL2/SDL.h>
import "C"

import (
	"encoding/binary"
	"image"
	"image/draw"
	"log"
	"unsafe"

	"github.com/mewkiz/pkg/imgutil"
)

// A Surface is a collection of pixels.
type Surface struct {
	// Width and height of the surface.
	Width, Height int
	// C surface pointer.
	s *C.SDL_Surface
}

// NewSurface returns a new surface of the specified dimensions.
//
// Note: The Free method of the surface should be called when finished using it.
func NewSurface(width, height int) (s *Surface, err error) {
	s = &Surface{
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
	s.s = C.SDL_CreateRGBSurface(0, C.int(width), C.int(height), 32, r, g, b, a)
	if s.s == nil {
		return nil, getError()
	}
	return s, nil
}

// LoadSurface loads the provided image and returns a surface containing its
// pixels.
//
// Note: The Free method of the surface should be called when finished using it.
func LoadSurface(imgPath string) (s *Surface, err error) {
	img, err := imgutil.ReadFile(imgPath)
	if err != nil {
		return nil, err
	}
	return GetSurface(img)
}

// GetSurface returns a surface containing the pixels of the provided image.
//
// Note: The Free method of the surface should be called when finished using it.
func GetSurface(img image.Image) (s *Surface, err error) {
	rect := img.Bounds()
	var pix []uint8
	switch i := img.(type) {
	case *image.NRGBA:
		pix = i.Pix
	case *image.RGBA:
		pix = i.Pix
	default:
		log.Printf("sdl.GetSurface: using draw.Draw fallback for image type %T.\n", img)
		fallback := image.NewNRGBA(rect)
		draw.Draw(fallback, rect, img, image.ZP, draw.Src)
		pix = fallback.Pix
	}
	s, err = NewSurface(rect.Dx(), rect.Dy())
	if err != nil {
		return nil, err
	}
	C.memcpy(s.s.pixels, unsafe.Pointer(&pix[0]), C.size_t(len(pix)))
	return s, nil
}

// Free frees the surface.
func (s *Surface) Free() {
	C.SDL_FreeSurface(s.s)
}

// Draw draws the entire src surface onto the dst surface starting at the
// destination point dp.
func (src *Surface) Draw(dst *Surface, dp image.Point) (err error) {
	dr := image.Rect(dp.X, dp.Y, dp.X+src.Width, dp.Y+src.Height)
	return src.DrawRect(dst, dr, image.ZP)
}

// DrawRect fills the destination rectangle dr of the dst surface with
// corresponding pixels from the src surface starting at the source point sp.
func (src *Surface) DrawRect(dst *Surface, dr image.Rectangle, sp image.Point) (err error) {
	sr := image.Rect(sp.X, sp.Y, sp.X+dr.Dx(), sp.Y+dr.Dy())
	srcRect := cRect(sr)
	dstRect := cRect(dr)
	if C.SDL_BlitSurface(src.s, srcRect, dst.s, dstRect) != 0 {
		return getError()
	}
	return nil
}
