#!/bin/bash

# Usage:
#   ./generate-mpd.sh output/playlist.mpd input/file1.mp4 input/file2.mp4 ...
#
# Output:
#   creates output/playlist.mpd
#   creates output/file1_dashinit.mp4
#   creates output/file2_dashinit.mp4

output_name="$1"
shift
input_names="$@"

# The segment duration is the amount of time in between i-frames.
# If you change this, update gop_size in `transcode-video.sh`.
segment_duration_ms=2000

MP4Box -dash "$segment_duration_ms" -frag "$segment_duration_ms" -rap -profile onDemand -out "$output_name" $input_names
