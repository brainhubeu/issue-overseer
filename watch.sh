#!/bin/sh

while true
do
  time ./issue-overseer $GITHUB_ORGANIZATION
  sleep 120
done
