package win

// #include <SDL2/SDL.h>
import "C"

import (
	"github.com/mewmew/we"
)

// goButton returns the corresponding we.Button for the provided SDL mouse
// button index.
func goButton(index C.Uint8) (button we.Button) {
	switch index {
	case C.SDL_BUTTON_LEFT:
		return we.ButtonLeft
	case C.SDL_BUTTON_RIGHT:
		return we.ButtonRight
	case C.SDL_BUTTON_MIDDLE:
		return we.ButtonMiddle
	case C.SDL_BUTTON_X1:
		return we.Button4
	case C.SDL_BUTTON_X2:
		return we.Button5
	}
	// All unknown SDL mouse buttons are mapped to buttons above we.Button5.
	return we.Button5 + we.Button(index-C.SDL_BUTTON_X2)
}

// goButtonFromState returns the corresponding we.Button for the provided SDL
// mouse button state bitfield.
//
// Note: Don't call this function with a state value of 0.
func goButtonFromState(state C.Uint32) (button we.Button) {
	if state&C.SDL_BUTTON_LMASK != 0 {
		return we.ButtonLeft
	}
	if state&C.SDL_BUTTON_RMASK != 0 {
		return we.ButtonRight
	}
	if state&C.SDL_BUTTON_MMASK != 0 {
		return we.ButtonMiddle
	}
	if state&C.SDL_BUTTON_X1MASK != 0 {
		return we.Button4
	}
	if state&C.SDL_BUTTON_X2MASK != 0 {
		return we.Button5
	}
	// All unknown SDL mouse buttons are mapped to buttons above we.Button5.
	for index := uint(6); index < 32; index++ {
		mask := C.Uint32(1 << (index - 1))
		if state&mask != 0 {
			return we.Button5 + we.Button(index-C.SDL_BUTTON_X2)
		}
	}
	panic("goButtonFromState: invalid mouse button state: 0.")
}
