#!/bin/bash

# Script to generate test.md5
#
# You shouldn't need to rerun this script unless you change the test files.

rm -f test.md5
for i in *.flac; do
	echo "Generating hash for $i..."
	sox $i -e signed $i.raw
	md5sum $i.raw >> test.md5
	rm -f $i.raw
done
