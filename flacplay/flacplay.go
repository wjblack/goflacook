// Test programs to use the mewkiz FLAC decoder and the gordonklaus portaudio
// library to convert/play FLAC files.

// GoFlaCook is (C)2017 by BJ Black, released under WTFPL.  See COPYING.

// flacplay plays one or more FLAC files
package main

import (
	"github.com/gordonklaus/portaudio"
	"github.com/mewkiz/flac"
	"github.com/wjblack/goflacook"
	"fmt"
	"os"
)

var output []int32
var outstream *portaudio.Stream

func main() {
	// Make sure we are doing at least one play and that there are
	// pairs of arguments (foo -> bar, baz -> bingo, etc)
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <file> [file...]\n",
			os.Args[0])
		os.Exit(-1)
	}

	outputter := goflacook.NewOutputter(proc)

	// Just run the flacplay func for each pair of filenames
	for i:=1; i<len(os.Args); i++ {
		flacplay(outputter, os.Args[i])
	}
}

// proc shoves samples out to the bufio writer.
func proc(stream *flac.Stream, samples []int32) error {
	output = samples
	err := outstream.Write()
	return err
}

// flacplay opens/decodes the FLACfile stated in infile and outputs RAW data
// to outfile.  Panics on any error.
func flacplay(outputter *goflacook.Outputter, filename string) {
	// Let the peanut gallery know what's up.
	fmt.Printf("Playing %s...\n", filename)

	// Open the stream.
	chk("while opening infile", outputter.Init(filename))
	channels := outputter.Stream.Info.NChannels
	bps := outputter.Stream.Info.BitsPerSample
	rate := outputter.Stream.Info.SampleRate
	fmt.Printf("%d channels, %d bits/sample, %d sample rate\n", channels,
		bps, rate)

	portaudio.Initialize()
	defer portaudio.Terminate()

	stream, err := portaudio.OpenDefaultStream(0, int(channels), float64(rate), len(output), &output)
	chk("opening portaudio", err)
	chk("starting stream", stream.Start())
	outstream = stream

	chk("playing sample", outputter.MainLoop())
}

// chk dies with a formatted error if err is not nil.
func chk(op string, err error) {
	if err != nil {
		panic(fmt.Sprintf("Error while %s: %s", op, err.Error()))
	}
}
