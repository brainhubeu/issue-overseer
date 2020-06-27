#!/bin/sh

while true
do
  time ./issue-overseer $GITHUB_ORGANIZATION
  sleep $SLEEP_IN_SECONDS
done
