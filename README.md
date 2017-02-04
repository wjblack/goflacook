GoFlaCook
=========
This is the Go FLAC Cookbook sample code.

It's a small cookbook of stuff that you can do with
[https://github.com/eaburns/flac] and
[https://github.com/gordonklaus/portaudio].

I implement a simple FLAC file player and a FLAC-to-RAW sample converter
utility, along with a test program that verifies that decoding matches exactly
the Xiph reference implementation (via precomputed MD5 hashes).


Code Status
-----------
[![Build Status](https://travis-ci.org/wjblack/goflacook.svg?branch=master)](https://travis-ci.org/wjblack/goflacook)


Installation
------------
The easiest way is to just:

`go get github.com/wjblack/goflacook/...`

If you want to test this (especially on a different arch), I definitely
recommend running "go test -v github.com/wjblack/goflacook".


Usage
-----
flacplay is pretty much just "flacplay file1 file2 file3..."

flac2raw takes in/out pairs, so "flac2raw infile1 outfile1 ..."

Note that flac2raw gives you a sox play() command line to replay the sound
if you want to test it out (plus you can steal the flags for sox conversion
commands, too).


Implementation
--------------
goflacook itself does most of the heavy lifting, you only need to:

1. Create a proc() function that corresponds to the Outputter signature
2. Instantiate NewOutputter(proc)
3. Run NewOutputter.Init(filename)
4. Run NewOutputter.MainLoop()
5. Clean up

Basically the test programs do exactly that.  If you want to create a new
library to extend that, go right ahead (this is WTFPL licensed, so do whatever
you want).


Notes
-----
I've definitely learned a bunch about the FLAC format and PortAudio.  The
big stuff:

1. Each FLAC frame has one subframe per channel.  So when outputting, you need
   to interleave yourself.  So spit out one sample (for however many bytes 1-3
   that is) from the first subframe, then one from the second, and so on.
2. The PortAudio binding behaves very oddly if you have a stream of 8-bit
   signed samples.  Better off converting those to []int16 or []int32 if
   possible.
3. [https://github.com/eaburns/flac] is somewhat nicer to use, but has a
   weird (presumably type conversion) bug when using 24-bit FLACs at time
   of writing that I was unable to debug.
