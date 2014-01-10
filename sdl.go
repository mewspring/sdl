// Package sdl provides simplified bindings for the SDL library version 2.0.
package sdl

// #cgo pkg-config: sdl2
// #include <SDL2/SDL.h>
import "C"

// InitFlag is a bitfield of subsystem initialization flags.
type InitFlag uint32

// SDL subsystem init flags.
const (
	InitAudio       = InitFlag(C.SDL_INIT_AUDIO)
	InitVideo       = InitFlag(C.SDL_INIT_VIDEO)
	InitNoParachute = InitFlag(C.SDL_INIT_NOPARACHUTE) // Don't catch fatal signals.
	InitEverything  = InitAudio | InitVideo
)

// Init initializes the subsystems specified by flags.
//
// Note: Init must be called before calling any other SDL function and Quit must
// be called when finished using the SDL library.
//
// Note: Unless the InitNoParachute flag is set, init will install cleanup
// signal handlers for some commonly ignored fatal signals (like SIGSEGV).
func Init(flags InitFlag) (err error) {
	if C.SDL_Init(C.Uint32(flags)) != 0 {
		return getError()
	}
	return nil
}

// Quit cleans up all initialized subsystems.
func Quit() {
	C.SDL_Quit()
}
