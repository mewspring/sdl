package texture

// #cgo pkg-config: sdl2
// #include <SDL2/SDL.h>
import "C"

import (
	"errors"
)

// getLastError returns the last error message.
func getLastError() error {
	return errors.New(C.GoString(C.SDL_GetError()))
}
