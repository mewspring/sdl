package audio

// #cgo pkg-config: SDL2_mixer
// #include <SDL2/SDL_mixer.h>
import "C"

// A Sound represents an active sound with a dedicated mixer channel.
type Sound struct {
	// End is a channel on which true is sent when the sound has finished and
	// reached the end.
	End <-chan bool
	end chan bool
	// The dedicated mixer channel of the sound.
	channel C.int
}

// Pause pauses the playback of the sound.
func (snd *Sound) Pause() {
	if !snd.isValid() {
		return
	}
	C.Mix_Pause(snd.channel)
}

// Resume resumes the playback of the sound.
func (snd *Sound) Resume() {
	if !snd.isValid() {
		return
	}
	C.Mix_Resume(snd.channel)
}

// Stop stops the playback of the sound and releases its dedicated mixer
// channel.
func (snd *Sound) Stop() {
	if !snd.isValid() {
		return
	}
	C.Mix_HaltChannel(snd.channel)
}
