package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
)

var (
	speakerInitialized bool
)

func playSound(filepath string) error {
	// Open the audio file
	f, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer f.Close()

	// Determine the file format and decode
	var streamer beep.StreamSeekCloser
	var format beep.Format

	if ext := filepath[len(filepath)-4:]; ext == ".mp3" {
		streamer, format, err = mp3.Decode(f)
	} else if ext == ".wav" {
		streamer, format, err = wav.Decode(f)
	} else {
		return fmt.Errorf("unsupported file format: %s", ext)
	}

	if err != nil {
		return fmt.Errorf("error decoding audio: %v", err)
	}
	defer streamer.Close()

	// Initialize the speaker if it hasn't been initialized yet
	if !speakerInitialized {
		err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		if err != nil {
			return fmt.Errorf("error initializing speaker: %v", err)
		}
		speakerInitialized = true
	}

	// Play the audio
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	// Wait for the audio to finish playing
	<-done

	return nil
}
