package main

import (
	"github.com/gordonklaus/portaudio"
	"github.com/ldeng7/go-fc/core"
)

type Audio struct {
	sampRate uint16
	stream   *portaudio.Stream
	source   *core.ApuDataQueue
}

func newAudio() (*Audio, error) {
	a := &Audio{}
	var err error
	defer func() {
		if err != nil {
			a.deInit()
		}
	}()

	portaudio.Initialize()
	hostApi, err := portaudio.DefaultHostApi()
	if err != nil {
		return nil, err
	}
	p := portaudio.HighLatencyParameters(nil, hostApi.DefaultOutputDevice)
	p.Output.Channels = 1
	a.sampRate = uint16(p.SampleRate)

	a.stream, err = portaudio.OpenStream(p, func(buf []int16) {
		a.source.Dequeue(buf)
	})
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *Audio) deInit() {
	if a.stream != nil {
		a.stream.Close()
	}
	portaudio.Terminate()
}
