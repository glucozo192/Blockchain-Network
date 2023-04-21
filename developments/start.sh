#!/bin/bash

export MAX_NODES=5
export SAMPLE_SIZE=4
export QUORUM_SIZE=3
export DECISION_THRESHOLD=5

rm -rf ./logs/*
# nodes=$(seq 1 $MAX_NODES)

for ((i = 0; i < MAX_NODES; i++)); do
  export PORT=$((9000 + $i))
  export NODE_ID=$((9000 + $i))
  go run ./cmd/srv/. &
done

wait
