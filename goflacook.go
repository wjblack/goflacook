// Test programs to use the mewkiz FLAC decoder and the gordonklaus portaudio
// library to convert/play FLAC files.

// GoFlaCook is (C)2017 by BJ Black, released under WTFPL.  See COPYING.

package goflacook

import (
	"github.com/mewkiz/flac"
	"github.com/mewkiz/flac/frame"
	"fmt"
	"io"
)

// Outputter is where the bits to be output live.  We assume three phases:
//   1. Initialization (e.g. opening an output file)
//   2. Frame processing (run once per FLAC frame)
//   3. Cleanup (e.g. flushing buffers, etc)
type Outputter struct {
	// Proc is the callback to run with each frame as it's being processed.
	Proc func(*flac.Stream, []int32) error
	// Stream is the FLAC stream being played.
	Stream *flac.Stream
	// interleaver specifies a function pointer to whichever sample
	// integrator we're using.
	interleaver func(*frame.Frame) []int32
}

// NewOutputter is a simple constructor for Outputter.
func NewOutputter(proc func(*flac.Stream, []int32) error) *Outputter {
	return &Outputter{ proc, nil, nil }
}

// interleave8 assumes 8 bits per sample and integrates all channels into a
// single stream.
func interleave8(frame *frame.Frame) []int32 {
	channels := len(frame.Subframes)
	ret := make([]int32, int(frame.BlockSize) * channels)
	for i:=0; i<int(frame.BlockSize); i++ {
		for j:=0; j<channels; j++ {
			sample := frame.Subframes[j].Samples[i]
			ret[i*channels+j] = sample << 24
		}
	}
	return ret
}

// interleave16 assumes 16 bits per sample and integrates all channels into a
// single stream.
func interleave16(frame *frame.Frame) []int32 {
	channels := len(frame.Subframes)
	ret := make([]int32, int(frame.BlockSize) * channels)
	for i:=0; i<int(frame.BlockSize); i++ {
		for j:=0; j<channels; j++ {
			sample := frame.Subframes[j].Samples[i]
			ret[i*channels+j] = sample << 16
		}
	}
	return ret
}

// interleave24 assumes 24 bits per sample and integrates all channels into a
// single stream.
func interleave24(frame *frame.Frame) []int32 {
	channels := len(frame.Subframes)
	ret := make([]int32, int(frame.BlockSize) * channels)
	for i:=0; i<int(frame.BlockSize); i++ {
		for j:=0; j<channels; j++ {
			sample := frame.Subframes[j].Samples[i]
			ret[i*channels+j] = sample << 8
		}
	}
	return ret
}

// Init is run first by the various clients.
func (outputter *Outputter) Init(filename string) error {
	// Open the FLAC file using the library's Open() func.
	stream, err := flac.Open(filename)
	if err != nil { return err }
	switch stream.Info.BitsPerSample {
		case 8:
			outputter.interleaver = interleave8
		case 16:
			outputter.interleaver = interleave16
		case 24:
			outputter.interleaver = interleave24
		default:
			panic(fmt.Sprintf("Don't know how to deal with a %d-" +
				"bit sample type.",
				stream.Info.BitsPerSample))
	}
	outputter.Stream = stream
	return nil
}

// MainLoop processes all the frames in the stream.
func (outputter *Outputter) MainLoop() error {
	// The main loop.  Suck in the next bit of data, short-circuit if
	// we've hit EOF and otherwise spit it out to the writer.
	for {
		// Get the next blob of samples.
		frame, err := outputter.Stream.ParseNext()
		if err != nil {
			if err == io.EOF { return nil }
			return err
		}
		// Note that we have to interleave the samples ourselves, as
		// each subframe represents a channel in the output.  Hopefully
		// the FLAC file has its channels arranged in a reasonable
		// order...
		err = outputter.Proc(outputter.Stream,
			outputter.interleaver(frame))
		if err != nil { return err }
	}
	return nil
}
