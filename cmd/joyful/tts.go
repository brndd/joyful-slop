package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"git.annabunches.net/annabunches/joyful/internal/logger"
	"github.com/ebitengine/oto/v3"
	flag "github.com/spf13/pflag"
)

type TTSOptions struct {
	Disabled bool
	Voice    string
	Volume   int
	Pitch    int
	Range    int
	Speed    int
}

type TTS struct {
	options *TTSOptions
	otoCtx  *oto.Context
	phrases map[string][]byte
}

const (
	playbackCheckIntervalMs = 100
)

// TODO: make most of this configurable via file
func addTTSFlags() *TTSOptions {
	ops := &TTSOptions{}

	flag.BoolVar(&ops.Disabled, "no-tts", false, "Disable text-to-speech.")
	flag.StringVar(&ops.Voice, "tts-voice", "en", "Which voice to use for TTS; see 'espeak --voices' for a full list of options.")
	flag.IntVar(&ops.Volume, "tts-volume", 100, "Text to speech volume")
	flag.IntVar(&ops.Pitch, "tts-pitch", 50, "Text to speech volume")
	flag.IntVar(&ops.Range, "tts-range", 50, "Text to speech volume")
	flag.IntVar(&ops.Range, "tts-speed", 175, "Text to speech speaking speed (in words per minute)")

	return ops
}

func makeOtoContext() (*oto.Context, error) {
	op := &oto.NewContextOptions{
		SampleRate:   22050,
		ChannelCount: 1,
		Format:       oto.FormatSignedInt16LE,
	}

	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		return nil, err
	}
	<-readyChan // wait for initialization

	return otoCtx, nil
}

func newTTS(ops *TTSOptions) (*TTS, error) {
	if ops.Disabled {
		return nil, nil
	}

	context, err := makeOtoContext()
	if err != nil {
		return nil, err
	}

	return &TTS{
		options: ops,
		otoCtx:  context,
		phrases: make(map[string][]byte),
	}, nil
}

func (t *TTS) AddMessage(msg string) {
	// TODO: need to get lots of input validation in here
	// We execute `espeak-ng` directly because extant libraries produce terrible output
	// compared to the command-line utility. This also gives us a chance to
	cmd := exec.Command(
		"espeak-ng", "--stdout",
		"-v", t.options.Voice,
		"-a", strconv.Itoa(t.options.Volume),
		"-p", strconv.Itoa(t.options.Pitch),
		"-P", strconv.Itoa(t.options.Range),
		"-s", strconv.Itoa(t.options.Speed),
		msg,
	)

	wavData, err := cmd.Output()
	if err != nil {
		logger.LogError(err, "Failed to create TTS data")
		return
	}

	t.phrases[msg] = wavData
}

// "Say" generates TTS audio and plays it in a go routine
func (t *TTS) Say(msg string) error {
	if _, ok := t.phrases[msg]; !ok {
		return fmt.Errorf("tried to play non-buffered phrase '%s'", msg)
	}

	go func(buf []byte) {
		buffer := bytes.NewBuffer(buf)
		player := t.otoCtx.NewPlayer(buffer)

		volume := 0.0
		player.SetVolume(volume)
		player.Play()

		// Gradually ramp up the volume to avoid harsh clicks
		for player.Volume() < 1.0 {
			volume += 0.01
			if volume > 1.0 {
				volume = 1.0
			}

			player.SetVolume(volume)
			time.Sleep(1 * time.Millisecond)
		}

		for player.IsPlaying() {
			time.Sleep(playbackCheckIntervalMs * time.Millisecond)
		}
	}(t.phrases[msg])

	return nil
}
