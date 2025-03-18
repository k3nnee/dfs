#!/bin/bash

function shut_down() {
  FOLDER=$(basename "$PWD")
  if [ "$FOLDER" != "naive-dfs" ]; then
    echo "Navigating into correct location ..."
    cd "$(dirname "$(realpath "$0")")" || exit 1
    echo "Navigated to $PWD"
  fi

  echo "Shutting down the containers ..."
  docker compose down || return 1

  if [ "$1" == "--clean-up" ]; then
    clean_up
  elif [ $"$1" ]; then
    read -r -p "Tag detected, did you mean --clean-up? (y/n): " confirm
    if [ "$confirm" == "y" ]; then
      clean_up
    fi
  fi

  echo "Complete"
}

function handle_error() {
  echo "Error occurred"
  exit 1
}

function clean_up() {
  echo "Cleaning up ..."
  docker system prune -a --volumes -f || return 1
}

shut_down "$@" || handle_error
