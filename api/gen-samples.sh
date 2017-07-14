#!/bin/bash

# This is really just an automation script using the follow ffmpeg command:
# ffmpeg -i input.mp4
#        -pix_fmt yuv420p
#        -movflags +faststart
#        -c:v libx264
#        -crf $quality
#        -vf "scale=w=$width:h=$height:force_original_aspect_ratio=decrease, scale=w=trunc(iw/2)*2:h=trunc(ih/2)*2"
#        -x264opts "keyint=$framerate:min-keyint=$framerate:no-scenecut"
#        -c:a copy
#        output.mp4


cd $1
input_name="$2.mp4"
output_name_pre="$2-"
output_name_post=".mp4"
framerate=30

# The GOP size is the number of frames between I-frames. This is defined as
# framerate * seconds between I-frames. The latter is defined in
# gen-mpeg-dash.sh as 2000ms (segment duration), so:
# 30 fps * 2 seconds = 60 frames
gop_size=60


function run_ffmpeg () {
	ffmpeg -i $input_name -pix_fmt yuv420p -movflags +faststart -c:v libx264 -r $framerate -crf $3 -vf "scale=w=$1:h=$2:force_original_aspect_ratio=decrease, scale=w=trunc(iw/2)*2:h=trunc(ih/2)*2" -x264opts "keyint=$gop_size:min-keyint=$gop_size:no-scenecut" -c:a copy "${output_name_pre}${2}-${3}${output_name_post}"
}

run_ffmpeg 480 360 17
run_ffmpeg 480 360 23
run_ffmpeg 480 360 29
run_ffmpeg 854 480 17
run_ffmpeg 854 480 23
run_ffmpeg 854 480 29
run_ffmpeg 1280 720 17
run_ffmpeg 1280 720 23
run_ffmpeg 1280 720 29
run_ffmpeg 1920 1080 17
run_ffmpeg 1920 1080 23
run_ffmpeg 1920 1080 29
