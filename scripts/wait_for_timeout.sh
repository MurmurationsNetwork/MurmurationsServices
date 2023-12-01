#!/bin/bash
# This script checks the HTTP status code of a given host until it matches a
# desired status code or until a specified timeout. It uses exponential backoff
# for retrying.

set -eux

declare HOST=$1
declare STATUS=$2
declare TIMEOUT=$3
declare SLEEP_TIME=2 # initial sleep time in seconds
declare MAX_SLEEP=120 # maximum sleep time in seconds

HOST=$HOST STATUS=$STATUS timeout --foreground -s TERM $TIMEOUT bash -c \
    'while [[ ${STATUS_RECEIVED} != ${STATUS} ]]; do
        sleep $SLEEP_TIME && \
        STATUS_RECEIVED=$(curl -s -o /dev/null -L -w ''%{http_code}'' ${HOST}) && \
        echo "received status: $STATUS_RECEIVED";
        if [[ ${STATUS_RECEIVED} != ${STATUS} ]]; then
            SLEEP_TIME=$((SLEEP_TIME * 2 > MAX_SLEEP ? MAX_SLEEP : SLEEP_TIME * 2))
        fi
     done;
     echo success with status: $STATUS_RECEIVED'
