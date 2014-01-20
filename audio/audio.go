// Package audio provides support for audio playback.
//
// The audio library is automatically initialized when imported. Therefore the
// Quit function must be called when finished using the library.
package audio

// #cgo pkg-config: sdl2 SDL2_mixer
// #include <SDL2/SDL.h>
// #include <SDL2/SDL_mixer.h>
//
// extern void onChannelFinished(int);
//
// static void callback(int channel) {
//    onChannelFinished(channel);
// }
//
// static void initCallback() {
//    Mix_ChannelFinished(callback);
// }
import "C"

import (
	"log"
)

func init() {
	err := initAudio()
	if err != nil {
		log.Fatalln(err)
	}
}

// initAudio initializes the audio subsystem.
//
// Note: The Quit function must be called when finished using the audio library.
func initAudio() (err error) {
	// Initialize the audio subsystem.
	if C.SDL_InitSubSystem(C.SDL_INIT_AUDIO) != 0 {
		return getSDLError()
	}

	// Open the audio output device with a frequency of 44.1 kHz and stereo
	// output channels (2).
	if C.Mix_OpenAudio(44100, C.MIX_DEFAULT_FORMAT, 2, 4096) != 0 {
		return getMixError()
	}

	// Initialize callback for channel finished playing events.
	C.initCallback()

	return nil
}

// Quit quits the audio subsystem.
//
// Note: The Quit function must be called when finished using the audio library.
func Quit() {
	// Close the audio output device.
	C.Mix_CloseAudio()

	// Quit all initialized file formats.
	for C.Mix_Init(0) != 0 {
		C.Mix_Quit()
	}

	// Quit the audio subsystem.
	C.SDL_QuitSubSystem(C.SDL_INIT_AUDIO)
}
