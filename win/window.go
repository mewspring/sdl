// Package win handles window creation, drawing and events. The window events
// are defined in a dedicated package located at:
//    github.com/mewmew/we
//
// The library uses a small subset of the features provided by SDL version 2.0.
// For the sake of simplicity support for multiple windows has intentionally
// been left out.
package win

// #cgo pkg-config: sdl2
// #include <stdlib.h>
// #include <string.h>
// #include <SDL2/SDL.h>
//
// static SDL_Rect * makeRectArray(int n) {
//    return calloc(n, sizeof(SDL_Rect));
// }
//
// static void setArrayRect(SDL_Rect *rects, SDL_Rect *rect, int i) {
//    void *dst = rects + i*sizeof(SDL_Rect);
//    memcpy(dst, rect, sizeof(SDL_Rect));
// }
import "C"

import (
	"image"
	"image/color"
)

// WindowFlag is a bitfield of window flags.
type WindowFlag uint32

// Window flags.
const (
	// Resizeable states that the window can be resized.
	Resizeable WindowFlag = C.SDL_WINDOW_RESIZABLE
	// FullScreen states that the window is in full screen mode.
	FullScreen WindowFlag = C.SDL_WINDOW_FULLSCREEN_DESKTOP
)

// win represents the graphics window which is opened through a call to Open. It
// is this single window that is utilized throughout the entire library.
var win *C.SDL_Window

// renderer represents the renderer of the window which is opened through a call
// to Open.
var renderer *C.SDL_Renderer

// Open opens a window with the specified dimensions and optional window flags.
// Only one window can be open at the same time. It is this single window that
// is utilized throughout the entire library. By default the window is not
// resizeable.
//
// Note: The Close function must be called when finished using the window.
func Open(width, height int, flags ...WindowFlag) (err error) {
	if win != nil {
		panic("win.Open: the window has already been opened.")
	}

	// Initialize the SDL video subsystem.
	if C.SDL_Init(C.SDL_INIT_VIDEO) != 0 {
		return getSDLError()
	}

	// Set window initialization flags.
	var cFlags C.Uint32
	for _, flag := range flags {
		cFlags |= C.Uint32(flag)
	}
	// TODO(u): enable vsync by default?
	// Always enable vsync.
	cFlags |= C.SDL_RENDERER_PRESENTVSYNC

	// Open the window.
	if C.SDL_CreateWindowAndRenderer(C.int(width), C.int(height), cFlags, &win, &renderer) != 0 {
		return getSDLError()
	}
	if win == nil || renderer == nil {
		return getSDLError()
	}

	err = Clear()
	if err != nil {
		return err
	}

	return nil
}

// Close closes the window.
func Close() {
	C.SDL_DestroyRenderer(renderer)
	renderer = nil
	C.SDL_DestroyWindow(win)
	win = nil
	C.SDL_Quit()
}

// SetTitle sets the title of the window.
func SetTitle(title string) {
	C.SDL_SetWindowTitle(win, C.CString(title))
}

// TODO(u): is the Size function needed?

//// Size returns the size of the window.
//func Size() (width, height int) {
//	var cWidth, cHeight C.int
//	C.SDL_GetWindowSize(win, &cWidth, &cHeight)
//	return int(cWidth), int(cHeight)
//}

// TODO(u): rethink Screen, return image.Image? Make use of
// SDL_RenderReadPixels.

//// Screen returns the image associated with the window.
//func Screen() (screen *Image, err error) {
//	screen = new(Image)
//	screen.s = C.SDL_GetWindowSurface(win)
//	if screen.s == nil {
//		return nil, getSDLError()
//	}
//	screen.Width = int(screen.s.w)
//	screen.Height = int(screen.s.h)
//	return screen, nil
//}

// Update displays window updates on the screen.
func Update() {
	C.SDL_RenderPresent(renderer)
}

// Clear clears the screen and fills it with the drawing color.
func Clear() (err error) {
	if C.SDL_RenderClear(renderer) != 0 {
		return getSDLError()
	}
	return nil
}

// SetDrawColor sets the drawing color of the window.
func SetDrawColor(c color.Color) (err error) {
	r, g, b, a := c.RGBA()
	if C.SDL_SetRenderDrawColor(renderer, C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)) != 0 {
		return getSDLError()
	}
	return nil
}

// Draw draws the entire src image onto the window starting at the destination
// point dp.
func Draw(dp image.Point, src *Image) (err error) {
	dr := image.Rect(dp.X, dp.Y, dp.X+src.Width, dp.Y+src.Height)
	return DrawRect(dr, src, image.ZP)
}

// DrawRect fills the destination rectangle dr of the window with corresponding
// pixels from the src image starting at the source point sp.
func DrawRect(dr image.Rectangle, src *Image, sp image.Point) (err error) {
	// Set the screen as the current rendering target.
	if C.SDL_SetRenderTarget(renderer, nil) != 0 {
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
