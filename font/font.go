// Package font handles text rendering based on the size, style and color of
// fonts.
package font

// #cgo pkg-config: SDL2_ttf
// #include <SDL2/SDL_ttf.h>
import "C"

import (
	"image/color"
	"log"

	"github.com/mewmew/sdl/win"
)

func init() {
	// Initializes the font library. The win.Close function calls SDL_Quit, which
	// also quits the font library.
	if C.TTF_Init() != 0 {
		log.Fatalln(getError())
	}
}

// A Font describes a particular size, style and color in which to render text.
type Font struct {
	// The maximum pixel height of all glyphs in the font.
	Height int
	// The color of the text drawn with font.
	c C.SDL_Color
	// C font pointer.
	f *C.TTF_Font
}

// Load loads the provided TTF font of the specified font size. The default
// color of a font is black.
//
// Note: The Free method of the font should be called when finished using it.
func Load(fontPath string, fontSize int) (f *Font, err error) {
	// Set default color to black.
	f = &Font{
		c: cColor(color.Black),
	}
	f.f = C.TTF_OpenFont(C.CString(fontPath), C.int(fontSize))
	if f.f == nil {
		return nil, getError()
	}
	f.Height = int(C.TTF_FontHeight(f.f))
	return f, nil
}

// Free frees the font.
func (f *Font) Free() {
	C.TTF_CloseFont(f.f)
}

// SetColor sets the color of the text drawn with font.
func (f *Font) SetColor(c color.Color) {
	f.c = cColor(c)
}

// Render renders an image of the provided text in the style, size and color of
// the font.
//
// Note: The Free method of the returned image should be called when finished
// using it.
func (f *Font) Render(text string) (img *win.Image, err error) {
	s := C.TTF_RenderUTF8_Blended(f.f, C.CString(text), f.c)
	if s == nil {
		return nil, getError()
	}
	img = winImage(s)
	return img, nil
}

// RenderWidth calculates the resulting image width of text if rendered using
// the font. The height is always f.Height. No actual rendering is performed.
func (f *Font) RenderWidth(text string) (width int, err error) {
	var w C.int
	if C.TTF_SizeUTF8(f.f, C.CString(text), &w, nil) != 0 {
		return 0, getError()
	}
	return int(w), nil
}
