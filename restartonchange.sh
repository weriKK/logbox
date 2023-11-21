#!/bin/sh

sigint_handler()
{
  kill $PID
  exit
}

trap sigint_handler INT

while true; do
  go build -o logbox .
  ./logbox &

  PID=$!
  inotifywait -e modify -e move -e create -e delete -e attrib -r `pwd`

  while test -d "/proc/$PID";  do
    kill $PID
    sleep 0.1
  done
done

