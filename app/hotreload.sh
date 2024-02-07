#!/bin/bash

# Default values
extensions="html,css,js,go,rs"
sleep_time=1
run_command=""
build_command=""

# Parse options
while getopts "e:s:r:b:" opt; do
  case $opt in
    e) extensions="$OPTARG";;
    s) sleep_time="$OPTARG";;
    r) run_command="$OPTARG";;
    b) build_command="$OPTARG";;
    \?) echo "Invalid option -$OPTARG" >&2
        exit 1
        ;;
  esac
done

echo "Hotreload settings:"
echo "Run command: '$run_command'"
echo "Extensions: '$extensions'"
echo "Build command: '$build_command'"
echo "Sleep time: '$sleep_time' seconds"
echo remember to use single quotes for the commands

if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS detected, use shasum
    hash_command="shasum"
else
    # Default to sha1sum (linux)
    hash_command="sha1sum"
fi

# Function to gracefully kill the previous process
cleanup() {
  local pid=""
  if [[ -f .hotreload ]]; then
    read -r _ pid < .hotreload
    if [[ -n "$pid" && $(ps -p $pid -o args=) == *"$run_command"* ]]; then
      kill "$pid"
      wait "$pid"
    fi
  fi
}

stop() {
  echo "cleaning up"
  cleanup
  exit 0
}

# Trap TERM signal and call cleanup function
trap stop TERM INT

if [[ -f .hotreload ]]; then
  rm .hotreload
fi

# Function to calculate cumulative hash of specified files
calculate_hash() {
  local current_hash=""
  IFS=',' read -ra ADDR <<< "$extensions"
  for ext in "${ADDR[@]}"; do
    for file in $(find . -type f -name "*.$ext" | sort); do
      if [[ -z "$current_hash" ]]; then
        current_hash=$(${hash_command} "$file" | awk '{print $1}')
      else
        current_hash=$(echo -n "$current_hash$(${hash_command} "$file" | awk '{print $1}')" | ${hash_command} | awk '{print $1}')
      fi
    done
  done
  echo "$current_hash"
}

# Main loop
while true; do
  current_hash=$(calculate_hash)

  if [[ -f .hotreload ]]; then
    read -r previous_hash previous_pid < .hotreload
  fi

  if [[ "$current_hash" != "$previous_hash" ]]; then
    echo "Hot-Reloading..."

    cleanup

    if [[ -n $build_command ]]; then
      $build_command &
      wait $!
    fi
    
    $run_command &
    command_pid=$!
    echo "$current_hash $command_pid" > .hotreload
    echo "started with pid $command_pid"
  fi

  sleep "$sleep_time"
done
