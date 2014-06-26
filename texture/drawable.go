package texture

// #cgo pkg-config: sdl2
// #include <SDL2/SDL.h>
import "C"

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"unsafe"

	"github.com/mewkiz/pkg/imgutil"
	"github.com/mewmew/wandi"
)

// Drawable represent a drawable texture. It implements the wandi.Drawable and
// wandi.Image interfaces.
type Drawable struct {
	// A drawable GPU texture.
	tex *C.SDL_Texture
}

// NewDrawable creates a drawable texture of the specified dimensions.
//
// Note: The Free method of the texture must be called when finished using it.
func NewDrawable(width, height int) (tex Drawable, err error) {
	tex.tex, err = create(width, height, true)
	if err != nil {
		return Drawable{}, err
	}
	return tex, nil
}

// LoadDrawable loads the provided file and converts it into a drawable texture.
//
// Note: The Free method of the texture must be called when finished using it.
func LoadDrawable(path string) (tex Drawable, err error) {
	src, err := imgutil.ReadFile(path)
	if err != nil {
		return Drawable{}, err
	}
	return ReadDrawable(src)
}

// ReadDrawable reads the provided image and converts it into a drawable
// texture.
//
// Note: The Free method of the texture must be called when finished using it.
func ReadDrawable(src image.Image) (tex Drawable, err error) {
	// Use fallback conversion for unknown image formats.
	rgba, ok := src.(*image.RGBA)
	if !ok {
		return ReadDrawable(fallback(src))
	}

	// Use fallback conversion for subimages.
	width, height := rgba.Rect.Dx(), rgba.Rect.Dy()
	if rgba.Stride != 4*width {
		return ReadDrawable(fallback(src))
	}

	// Create a drawable texture based on the pixels of the src image.
	tex, err = NewDrawable(width, height)
	if err != nil {
		return Drawable{}, err
	}
	pix := unsafe.Pointer(&rgba.Pix[0])
	if C.SDL_UpdateTexture(tex.tex, nil, pix, C.int(rgba.Stride)) != 0 {
		tex.Free()
		return Drawable{}, fmt.Errorf("texture.ReadDrawable: %v", getLastError())
	}

	return tex, nil
}

// Free frees the texture.
func (tex Drawable) Free() {
	C.SDL_DestroyTexture(tex.tex)
}

// Width returns the width of the texture.
func (tex Drawable) Width() int {
	var width C.int
	if C.SDL_QueryTexture(tex.tex, nil, nil, &width, nil) != 0 {
		log.Fatalln("Image.Width: unable to locate texture width;", getLastError())
	}
	return int(width)
}

// Height returns the height of the texture.
func (tex Drawable) Height() int {
	var height C.int
	if C.SDL_QueryTexture(tex.tex, nil, nil, nil, &height) != 0 {
		log.Fatalln("Image.Height: unable to locate texture height;", getLastError())
	}
	return int(height)
}

// Draw draws the entire src image onto the dst texture starting at the
// destination point dp.
func (dst Drawable) Draw(dp image.Point, src wandi.Image) (err error) {
	sr := image.Rect(0, 0, src.Width(), src.Height())
	return dst.DrawRect(dp, src, sr)
}

// drawRect draws a subset of the src texture, as defined by the source
// rectangle sr, onto the dst texture starting at the destination point dp.
func drawRect(dst *C.SDL_Texture, dp image.Point, src *C.SDL_Texture, sr image.Rectangle) (err error) {
	ren, err := getRenderer()
	if err != nil {
		return err
	}
	if C.SDL_SetRenderTarget(ren, dst) != 0 {
		return fmt.Errorf("texture.drawRect: %v", getLastError())
	}
	defer C.SDL_SetRenderTarget(ren, nil)
	width, height := C.int(sr.Dx()), C.int(sr.Dy())
	srcrect := &C.SDL_Rect{
		x: C.int(sr.Min.X),
		y: C.int(sr.Min.Y),
		w: width,
		h: height,
	}
	dstrect := &C.SDL_Rect{
		x: C.int(dp.X),
		y: C.int(dp.Y),
		w: width,
		h: height,
	}
	if C.SDL_RenderCopy(ren, src, srcrect, dstrect) != 0 {
		return fmt.Errorf("texture.drawRect: %v", getLastError())
	}
	return nil
}

// DrawRect draws a subset of the src image, as defined by the source rectangle
// sr, onto the dst texture starting at the destination point dp.
func (dst Drawable) DrawRect(dp image.Point, src wandi.Image, sr image.Rectangle) (err error) {
	switch srcImg := src.(type) {
	case Drawable:
		return drawRect(dst.tex, dp, srcImg.tex, sr)
	case Image:
		return drawRect(dst.tex, dp, srcImg.tex, sr)
	default:
		return fmt.Errorf("Drawable.DrawRect: source type %T not yet supported", src)
	}
}

// Fill fills the entire texture with the provided color.
func (dst Drawable) Fill(c color.Color) (err error) {
	ren, err := getRenderer()
	if err != nil {
		return err
	}
	if C.SDL_SetRenderTarget(ren, dst.tex) != 0 {
		return fmt.Errorf("Drawable.Fill: %v", getLastError())
	}
	defer C.SDL_SetRenderTarget(ren, nil)
	r, g, b, a := c.RGBA()
	if C.SDL_SetRenderDrawColor(ren, C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)) != 0 {
		return fmt.Errorf("Drawable.Fill: %v", getLastError())
	}
	if C.SDL_RenderClear(ren) != 0 {
		return fmt.Errorf("Drawable.Fill: %v", getLastError())
	}
	return nil
}

// Image returns an image.Image representation of the texture.
func (tex Drawable) Image() (img image.Image, err error) {
	// TODO(u): Why is the resulting image upside-down?
	ren, err := getRenderer()
	if err != nil {
		return nil, err
	}
	if C.SDL_SetRenderTarget(ren, tex.tex) != 0 {
		return nil, fmt.Errorf("Drawable.Image: %v", getLastError())
	}
	defer C.SDL_SetRenderTarget(ren, nil)
	rect := image.Rect(0, 0, tex.Width(), tex.Height())
	dst := image.NewRGBA(rect)
	format := C.Uint32(C.SDL_PIXELFORMAT_ABGR8888)
	pix := unsafe.Pointer(&dst.Pix[0])
	if C.SDL_RenderReadPixels(ren, nil, format, pix, C.int(dst.Stride)) != 0 {
		return nil, fmt.Errorf("Drawable.Image: %v", getLastError())
	}
	return dst, nil
}
