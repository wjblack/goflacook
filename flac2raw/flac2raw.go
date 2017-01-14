// Test programs to use the mewkiz FLAC decoder and the gordonklaus portaudio
// library to convert/play FLAC files.

// GoFlaCook is (C)2017 by BJ Black, released under WTFPL.  See COPYING.

// flac2raw converts FLAC files to little-endian signed RAW sample files.
package main

import (
	"github.com/mewkiz/flac"
	"github.com/wjblack/goflacook"
	"bufio"
	"fmt"
	"os"
)

// writer is the stream we're writing to.
var writer *bufio.Writer

func main() {
	// Make sure we are doing at least one conversion and that there are
	// pairs of arguments (foo -> bar, baz -> bingo, etc)
	if len(os.Args) < 3 || len(os.Args) % 2 != 1 {
		fmt.Printf("Usage: %s <infile> <outfile> [infile outfile...]\n",
			os.Args[0])
		os.Exit(-1)
	}

	outputter := goflacook.NewOutputter(proc)

	// Just run the flac2raw func for each pair of filenames
	for i:=1; i<len(os.Args); i+=2 {
		flac2raw(outputter, os.Args[i], os.Args[i+1])
	}
}

// proc shoves samples out to the bufio writer.
func proc(stream *flac.Stream, samples []int32) error {
	for _, sample := range samples {
		err := writeSample(sample, stream.Info.BitsPerSample)
		if err != nil { return err }
	}
	return nil
}

// flac2raw opens/decodes the FLACfile stated in infile and outputs RAW data
// to outfile.  Panics on any error.
func flac2raw(outputter *goflacook.Outputter, infile, outfile string) {
	// Let the peanut gallery know what's up.
	fmt.Printf("Converting %s to %s...\n", infile, outfile)

	// Open the stream.
	chk("while opening infile", outputter.Init(infile))
	channels := outputter.Stream.Info.NChannels
	bps := outputter.Stream.Info.BitsPerSample
	rate := outputter.Stream.Info.SampleRate
	fmt.Printf("%d channels, %d bits/sample, %d sample rate\n", channels,
		bps, rate)

	// Open the output file.
	out, err := os.OpenFile(outfile,
		os.O_WRONLY | os.O_CREATE | os.O_EXCL, 0644)
	chk("opening outfile", err); defer out.Close()
	writer = bufio.NewWriter(out)

	chk("converting sample", outputter.MainLoop())

	// Don't forget we're using a bufio writer!
	writer.Flush()

	// Let the user know this one's done and how to replay it.
	fmt.Printf("Done!  If you have sox installed, play with:\n")
	fmt.Printf("    play -c %d -b %d -r %d -e signed %s\n", channels,
		bps, rate, outfile)
}

// writeSample writes a single sample in the proper (signed, 1/2/3-byte)
// little-endian format.  We make the assumption that we can get away with
// one sample at a time to the writer because bufio is buffering for us.
func writeSample(sample int32, bps uint8) error {
	bytes := make([]byte,int(bps/8))
	switch bps {
		case 8:
			bytes[0] = byte((sample >> 24) & 0xff)
		case 16:
			bytes[0] = byte((sample >> 16) & 0xff)
			bytes[1] = byte((sample >> 24) & 0xff)
		case 24:
			bytes[0] = byte((sample >> 8) & 0xff)
			bytes[1] = byte((sample >> 16) & 0xff)
			bytes[2] = byte((sample >> 24) & 0xff)
	}
	_, err := writer.Write(bytes)
	return err
}

// chk dies with a formatted error if err is not nil.
func chk(op string, err error) {
	if err != nil {
		panic(fmt.Sprintf("Error while %s: %s", op, err.Error()))
	}
}
