package audio

// #cgo pkg-config: sdl2 SDL2_mixer
// #include <SDL2/SDL.h>
// #include <SDL2/SDL_mixer.h>
import "C"

import (
	"io"
)

// A Stream is a sequence of audio samples.
type Stream struct {
	// C mixer chunk pointer.
	c *C.Mix_Chunk
}

// Open returns a new audio stream which reads and decodes its samples from the
// provided file.
//
// Note: The Close method should be called when done using the audio stream.
func Open(filePath string) (stream *Stream, err error) {
	stream = new(Stream)
	src := C.SDL_RWFromFile(C.CString(filePath), C.CString("rb"))
	if src == nil {
		return nil, getSDLError()
	}
	stream.c = C.Mix_LoadWAV_RW(src, 1)
	if stream.c == nil {
		return nil, getMixError()
	}
	return stream, nil
}

// New returns a new audio stream which reads and decodes its samples from the
// provided io.Reader.
//
// Note: The Close method should be called when done using the audio stream.
func New(r io.Reader) (stream *Stream, err error) {
	// TODO(u): implement using Mix_LoadWAV_RW.
	panic("audio.New: not yet implemented.")
}

// Close closes the audio stream.
func (stream *Stream) Close() {
	C.Mix_FreeChunk(stream.c)
}

// Play starts to play the audio stream in a dedicated channel, and returns a
// handle to the active sound.
func (stream *Stream) Play() (snd *Sound, err error) {
	// TODO(u): dynamically allocate mixer channels when needed using
	// Mix_AllocateChannels.
	//    - First figure out if there are any available channels.
	//    - If not, allocate new ones; possibly in powers of 2.
	//    - Remember to set the volume of the newly allocated channels.

	// Play the audio stream on the first available mixer channel once.
	channel := C.Mix_PlayChannelTimed(-1, stream.c, 0, -1)
	if channel == -1 {
		return nil, getMixError()
	}
	snd = &Sound{
		channel: channel,
	}
	// Create a buffered channel for the end event so we never have to wait for a
	// channel receive.
	snd.end = make(chan bool, 1)
	snd.End = snd.end
	active[channel] = snd
	return snd, nil
}

// active is a map from active channels to the sound of the channel.
var active = make(map[C.int]*Sound)

// isValid returns true if the sound has a dedicated channel, and false
// otherwise.
func (snd *Sound) isValid() bool {
	_, ok := active[snd.channel]
	return ok
}

//export onChannelFinished
func onChannelFinished(channel C.int) {
	// Remove the finished channel from the active channels map.
	snd, ok := active[channel]
	if ok {
		delete(active, channel)
		snd.end <- true
		close(snd.end)
	}
}
