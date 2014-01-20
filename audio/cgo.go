package audio

// #cgo pkg-config: sdl2 SDL2_mixer
// #include <SDL2/SDL.h>
// #include <SDL2/SDL_mixer.h>
import "C"

import (
	"errors"
)

// getSDLError returns the last SDL error message.
func getSDLError() (err error) {
	return errors.New(C.GoString(C.SDL_GetError()))
}

// getMixError returns the last mixer error message.
func getMixError() (err error) {
	return errors.New(C.GoString(C.Mix_GetError()))
}
