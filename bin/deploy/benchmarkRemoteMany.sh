#!/bin/bash

# With this file, multiple block sizes can be run consecutively.

ALGORITHM=$1
BUILD=$2
NEW_AWS=1
SSH_KEY="/local/home/yvonlanthen/.ssh/aws_global"

# Example usage:
# bash benchmarkRemoteMany.sh icc 1

chmod +777 benchmarkRemote.sh

for (( j=1000; j<=5000*1000; j+=1000*1000))
do
    echo $ALGORITHM ${j} is launched!

    cat config.json | jq '.payload_size = ($v|tonumber)' --arg v ${j} | sponge config.json

    ./benchmarkRemote.sh $ALGORITHM $j 0 $BUILD

done

