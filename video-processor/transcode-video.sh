#!/bin/bash

# Usage:
#   ./transcode-video.sh path/to/input.mp4
#
# Output:
#   creates path/to/input_video_360_28.mp4
#   creates path/to/input_video_480_23.mp4
#   creates path/to/input_video_720_20.mp4
#   creates path/to/input_video_1080_17.mp4
#   creates path/to/input_audio_192.m4a
#
# This is really just an automation script using the follow ffmpeg command:
# ffmpeg -i input.mp4
#        -pix_fmt yuv420p
#        -movflags +faststart
#        -c:v libx264
#        -crf $quality
#        -preset slow
#        -vf "scale=w=$width:h=$height:force_original_aspect_ratio=decrease, scale=w=trunc(iw/2)*2:h=trunc(ih/2)*2"
#        -x264opts "keyint=$framerate:min-keyint=$framerate:no-scenecut"
#        -an
#        output.mp4

input_name="$1"
output_name_video_pre="${1%.*}_video"
output_name_video_post=".mp4"
output_name_audio_pre="${1%.*}_audio"
output_name_audio_post=".m4a"
framerate=60

# The GOP size is the number of frames between i-frames. This is defined as
# framerate * seconds between i-frames. The latter is defined in
# transcode-video.sh as 2000ms (segment duration), so:
# 60 fps * 2 seconds = 120 frames
gop_size=120

function ffmpeg_video () {
	ffmpeg -i "$input_name" -preset slow -an -c:v libx264 -r "$framerate" -crf "$3" -vf "scale=w=$1:h=$2:force_original_aspect_ratio=decrease, scale=w=trunc(iw/2)*2:h=trunc(ih/2)*2" -x264opts "keyint=$gop_size:min-keyint=$gop_size:no-scenecut" -pix_fmt yuv420p "${output_name_video_pre}_${2}_${3}${output_name_video_post}"
}

function ffmpeg_audio () {
	ffmpeg -i "$input_name" -vn -c:a aac -b:a "${1}k" "${output_name_audio_pre}_${1}${output_name_audio_post}"
}

ffmpeg_video 480 360 28
ffmpeg_video 854 480 23
ffmpeg_video 1280 720 20
ffmpeg_video 1920 1080 18
ffmpeg_audio 192
