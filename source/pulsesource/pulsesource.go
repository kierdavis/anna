// PulseAudio audio source
package pulsesource

import (
	"github.com/mesilliac/pulse-simple"
)

type Source struct {
	stream *pulse.Stream
}

func New(sampleRate float64) (src Source, err error) {
	ss := pulse.SampleSpec{
		Format: pulse.SAMPLE_U8,
		Rate: uint32(sampleRate),
		Channels: 2,
	}
	
	s, err := pulse.Capture("anna", "stereo input", &ss)
	if err != nil {
		return Source{}, err
	}
	
	return Source{s}, nil
}

func (src Source) Read(left []float64, right []float64) (err error) {
	if len(left) != len(right) {
		panic("left and right must have the same length")
	}
	
	buf := make([]byte, len(left)*2)
	_, err = src.stream.Read(buf)
	if err != nil {
		return err
	}
	
	j := 0
	for i := range left {
		l, r := buf[j], buf[j+1]
		left[i] = float64(l - 128) / 128.0
		right[i] = float64(r - 128) / 128.0
		j += 2
	}
	
	return nil
}

func (src Source) Close() {
	src.stream.Free()
}
