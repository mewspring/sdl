// TODO(u): rethink the API. Make it work more like win. audio.Play could be a
// convenience function around an OutputSystem/SpeakerSetup/Speaker/Output/...
// convenience function around an OutputSystem/SpeakerSetup/Speaker/Output/
// SpeakerSet...
//
// This would make it possible to take advantage of multiple speakers if they
// are present. It would also make it at least theoretically possible to play
// some sounds through one set of speakers and other sounds through another set
// of speakers.

// Package audio provides support for audio playback.
//
// NOTE: This is a work in progress. It is in the API design phase and no
// functionality has yet been implemented.
package audio

import (
	"io"
	"time"
)

// A Stream is a sequence of audio samples.
type Stream struct {
}

// New returns a new audio stream which reads and decodes its samples from the
// provided io.Reader.
func New(r io.Reader) (s *Stream, err error) {
	panic("audio.New: not yet implemented.")
}

// Open returns a new audio stream which reads and decodes its samples from the
// provided file.
//
// Note: The Close method should be called when done using the audio stream.
func Open(filePath string) (s *Stream, err error) {
	panic("audio.Open: not yet implemented.")
}

// TODO(u): document if Close is required to be called for client of both New
// and Open.

// Close closes the audio stream.
func (s *Stream) Close() {
	panic("Stream.Close: not yet implemented.")
}

// Play starts to play the audio stream from the current audio position.
func (s *Stream) Play() {
	panic("Stream.Play: not yet implemented.")
}

// Pause pauses the playback of the audio stream. The current audio position is
// stored.
func (s *Stream) Pause() {
	panic("Stream.Pause: not yet implemented.")
}

// Stop stops the playback of the audio stream. The current audio position is
// set to 0.
func (s *Stream) Stop() {
	panic("Stream.Stop: not yet implemented.")
}

// Seek seeks to the provided audio position.
func (s *Stream) Seek(pos time.Duration) (err error) {
	panic("Stream.Seek: not yet implemented.")
}

// Pos returns the current audio position of the audio stream.
func Pos() (pos time.Duration) {
	panic("Stream.Pos: not yet implemented.")
}
