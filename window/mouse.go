package window

// #include <SDL2/SDL.h>
import "C"

import (
	"log"

	"github.com/mewmew/we"
)

// weButton returns the corresponding we.Button for the provided SDL mouse
// button index.
func weButton(index C.Uint8) (button we.Button) {
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

	// Unknown mouse button.
	log.Printf("window.weButton: unknown mouse button %d.\n", button)
	return 0
}

// weButtonFromState returns the corresponding we.Button for the provided SDL
// mouse button state bitfield.
//
// Note: Don't call this function with a state value of 0.
func weButtonFromState(state C.Uint32) (button we.Button) {
	if state == 0 {
		panic("weButtonFromState: invalid mouse button state: 0.")
	}
	if state&C.SDL_BUTTON_LMASK != 0 {
		button |= we.ButtonLeft
	}
	if state&C.SDL_BUTTON_RMASK != 0 {
		button |= we.ButtonRight
	}
	if state&C.SDL_BUTTON_MMASK != 0 {
		button |= we.ButtonMiddle
	}
	if state&C.SDL_BUTTON_X1MASK != 0 {
		button |= we.Button4
	}
	if state&C.SDL_BUTTON_X2MASK != 0 {
		button |= we.Button5
	}
	return button
}
