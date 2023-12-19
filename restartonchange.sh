#!/bin/sh

sigint_handler()
{
  kill $PID
  exit
}

trap sigint_handler INT

while true; do
  printf "\n\n"
  
  go build -o bin/logbox .
  ./bin/logbox &

  PID=$!
  inotifywait -e modify -e move -e create -e delete -e attrib --exclude "logs.*\.db" -r `pwd`

  while test -d "/proc/$PID";  do
    kill $PID
    sleep 0.1
  done
done

