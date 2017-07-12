#!/bin/bash

output_name_pre="sample"
input_name_pre="sample-"
segment_name_pre="sample-segment_"

# If you change this, update gop_size in gen-samples.sh.
segment_duration_ms=2000


MP4Box -dash $segment_duration_ms -frag $segment_duration_ms -rap -profile onDemand $output_name_pre ${input_name_pre}*
