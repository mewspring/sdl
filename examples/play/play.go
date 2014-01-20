// play demonstrates how to open audio files and play sounds.
package main

import (
	"log"

	"github.com/mewkiz/pkg/goutil"
	"github.com/mewmew/sdl/audio"
)

func main() {
	// The audio library is automatically initialized when imported. Quit the
	// audio library on return.
	defer audio.Quit()

	err := play()
	if err != nil {
		log.Fatalln(err)
	}
}

// play demonstrates how to open audio files and play sounds.
func play() (err error) {
	dataDir, err := goutil.SrcDir("github.com/mewmew/sdl/examples/play/data")
	if err != nil {
		return err
	}

	// Load the sound file.
	s, err := audio.Open(dataDir + "/birds.ogg")
	if err != nil {
		return err
	}

	// Play the sound file.
	snd, err := s.Play()
	if err != nil {
		return err
	}

	// Wait until the sound has reached the end.
	<-snd.End

	return nil
}
