package win

// #cgo pkg-config: sdl2
// #include <SDL2/SDL.h>
import "C"

import (
	"errors"
	"image"
	"log"
	"unsafe"
)

// cRect converts a Go image.Rectangle to a C SDL_Rect.
func cRect(rect image.Rectangle) (cRect *C.SDL_Rect) {
	if rect == image.ZR {
		return nil
	}
	cRect = new(C.SDL_Rect)
	cRect.x = C.int(rect.Min.X)
	cRect.y = C.int(rect.Min.Y)
	cRect.w = C.int(rect.Max.X - rect.Min.X)
	cRect.h = C.int(rect.Max.Y - rect.Min.Y)
	return cRect
}

// getError returns the last error message.
func getError() (err error) {
	return errors.New(C.GoString(C.SDL_GetError()))
}

// nativeBigEndian is set to true if the native byte order of the system is big
// endian, and false otherwise.
var nativeBigEndian bool

// initNativeByteOrder determintes the native byte order of the system.
func initNativeByteOrder() (err error) {
	i := int32(0x01020304)
	p := (*byte)(unsafe.Pointer(&i))
	switch *p {
	case 0x01:
		nativeBigEndian = true
		return nil
	case 0x04:
		nativeBigEndian = false
		return nil
	}
	return errors.New("win.initNativeByteOrder: unable to determine native byte order")
}

func init() {
	err := initNativeByteOrder()
	if err != nil {
		log.Fatalln(err)
	}
}
