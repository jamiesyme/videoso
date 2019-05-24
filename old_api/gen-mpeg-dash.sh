#!/bin/bash

cd $1
output_name_pre="$2"
input_name_pre="$2-"

# If you change this, update gop_size in gen-samples.sh.
segment_duration_ms=2000


MP4Box -dash $segment_duration_ms -frag $segment_duration_ms -rap -profile onDemand $output_name_pre ${input_name_pre}*
