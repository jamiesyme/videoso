#!/bin/bash

# Usage:
#   ./transcode-image.sh path/to/input.jpg
#
# Output:
#   creates path/to/input_240.jpg
#   creates path/to/input_360.jpg
#   creates path/to/input_480.jpg
#   creates path/to/input_720.jpg
#   creates path/to/input_1080.jpg
#
# This is really just an automation script using the follow ImageMagick command:
# convert input.jpg
#         -strip
#         -thumbnail ${width}x${height}
#         output.jpg

input_name="$1"
output_name_pre="${1%.*}"
output_name_post=".jpg"

function convert_image () {
	convert "$input_name" -strip -thumbnail "${1}x${2}" "${output_name_pre}_${2}${output_name_post}"
}

convert_image 426 240
convert_image 640 360
convert_image 854 480
convert_image 1280 720
convert_image 1920 1080
