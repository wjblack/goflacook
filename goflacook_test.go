// Test of the goflacook package.  Basically we iterate through the entries in
// test.md5 and double-check that decoding the FLAC file to RAW samples comes
// up with the same output as the reference FLAC decoder (MD5s precomputed
// using mkhash.sh using the Xiph reference FLAC program).

package goflacook

import (
	"github.com/mewkiz/flac"
	"bufio"
	"crypto/md5"
	"fmt"
	"hash"
	"os"
	"regexp"
	"testing"
)

// md holds our MD5 hasher.
var md hash.Hash

// TestSamples tries to enumerate the files in test.md5 and verify that they
// decode/hash the same as the precomputed hashes.
func TestSamples(t *testing.T) {
	// Open the hashfile
	infile, err := os.Open("test.md5")
	if err != nil { t.Error(err.Error()); return }
	// Create a regex to extract filename and precomputed hash
	re, err := regexp.Compile("^(\\S+)\\s+(\\S+)\\.raw$")
	if err != nil { t.Error(err.Error()); return }
	// Run through each line of the file and execute the subtest
	// Lines that don't match expectations are skipped
	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		matches := re.FindStringSubmatch(scanner.Text())
		if len(matches) == 3 {
			HashAudio(t, matches[2], matches[1])
		}
	}
}

// HashAudio runs a subtest against a given FLAC file.
func HashAudio(t *testing.T, filename, md5sum string) {
	md = md5.New()
	outputter := NewOutputter(proc)
	err := outputter.Init(filename)
	if err != nil { t.Error(err.Error()); return }
	err = outputter.MainLoop()
	if err != nil { t.Error(err.Error()); return }
	hash := fmt.Sprintf("%x", md.Sum(nil))
	if hash == md5sum {
		t.Logf("Test on %s passed!", filename)
	} else {
		t.Errorf("%s: Expected %s, got %s", filename, md5sum, hash)
	}
}

// proc, required by Outputter, spits all samples into the hasher.
func proc(stream *flac.Stream, samples []int32) error {
	for _, sample := range samples {
		addSample(sample, stream.Info.BitsPerSample)
	}
	return nil
}

// addSample converts a sample to a []byte and shoves it into the hasher.
func addSample(sample int32, bps uint8) {
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
	md.Write(bytes)
}
