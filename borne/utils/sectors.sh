#!/bin/bash
 
set -euo pipefail
 
# This script helps during the build process to determine the start/end sectors for a partition
# or find the size of a partition in sectors
# It can run against a block device or raw image since parted is capable of this
 
# Usage:    ./sectors.sh [start/end/size] [DEVICE] [PARTNUM]
# Example:  ./sectors.sh start /dev/sda 2
 
USAGE="Usage: $0 [start/end/size] [DEVICE] [PARTNUM]"
 
STAT=${1:?$USAGE}
DEVICE=${2:?$USAGE}
PARTNUM=${3:?$USAGE}
 
case $STAT in
    start)  FIELD=2;;
    end)    FIELD=3;;
    size)   FIELD=4;;
    *)      
        echo "$USAGE" >&2
        exit 1
    ;;
esac
 
# The 'machine friendly' output from parted is not particularly nice to parse
# Select a specific line from the output based on the partition number
# Take the nth field depending on what stat was requested
 
SECTORS=$(parted --machine -s "$DEVICE" unit S print | awk NR==$((PARTNUM + 2)) | cut -f $FIELD -d ':')
 
# Echo with the 's' truncated
echo ${SECTORS/s}

