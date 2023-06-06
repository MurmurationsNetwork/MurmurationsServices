#!/bin/bash
# This is a bash script that checks the HTTP status code of a given host until
# it matches a desired status code or until a specified timeout.

set -eux

declare HOST=$1
declare STATUS=$2
declare TIMEOUT=$3

HOST=$HOST STATUS=$STATUS timeout --foreground -s TERM $TIMEOUT bash -c \
    'while [[ ${STATUS_RECEIVED} != ${STATUS} ]];\
        do sleep 30 && \
        STATUS_RECEIVED=$(curl -s -o /dev/null -L -w ''%{http_code}'' ${HOST}) && \
        echo "received status: $STATUS_RECEIVED"; \
    done;
    echo success with status: $STATUS_RECEIVED'
