package font

// #include <SDL2/SDL.h>
//
// SDL_Color colorAt(SDL_Color *colors, int i) {
//    return colors[i];
// }
import "C"

import (
	"image"
	"image/color"
	"log"
	"reflect"
	"unsafe"
)

// castImage returns the image.Image equivalent of the provided SDL surface. The
// returned image shares pixels with the original surface.
func castImage(s *C.SDL_Surface) (img image.Image) {
	width, height := int(s.w), int(s.h)
	stride := int(s.pitch)
	n := stride * height
	sh := reflect.SliceHeader{
		Data: uintptr(s.pixels),
		Len:  n,
		Cap:  n,
	}
	pix := *(*[]uint8)(unsafe.Pointer(&sh))
	rect := image.Rect(0, 0, width, height)
	switch s.format.format {
	case C.SDL_PIXELFORMAT_INDEX8:
		pal := goPalette(s.format.palette)
		img = &image.Paletted{
			Pix:     pix,
			Stride:  stride,
			Rect:    rect,
			Palette: pal,
		}
	case C.SDL_PIXELFORMAT_ARGB8888:
		img = &image.NRGBA{
			Pix:    pix,
			Stride: stride,
			Rect:   rect,
		}
	default:
		log.Fatalf("font.castImage: image format %d not yet supported.\n", s.format.format)
	}
	return img
}

// goPalette converts the provided C SDL_Palette to a Go color.Palette.
func goPalette(p *C.SDL_Palette) (pal color.Palette) {
	n := int(p.ncolors)
	pal = make([]color.Color, n)
	for i := 0; i < n; i++ {
		c := goColor(C.colorAt(p.colors, C.int(i)))
		pal[i] = c
	}
	return pal
}

// goColor converts the provided C SDL_Color to a Go color.Color.
func goColor(c C.SDL_Color) color.Color {
	return color.RGBA{uint8(c.r), uint8(c.g), uint8(c.b), uint8(c.a)}
}
