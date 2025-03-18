#!/bin/bash

function start_up() {
  FOLDER=$(basename "$PWD")
  if [ "$FOLDER" != "naive-dfs" ]; then
    echo "Navigating into correct location ..."
    cd "$(dirname "$(realpath "$0")")" || exit 1
    echo "Navigated to $PWD"
  fi

  docker compose up --build --force-recreate
}

start_up