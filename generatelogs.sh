#!/bin/sh

while true; do

    echo "Lorem ipsum dolor sit amet" | nc -q 0 localhost 8888
    sleep 1

done

