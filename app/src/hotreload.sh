#!/bin/bash

# Default values
file_grep="^.+\.(go|css|html|rs|js|templ)$"
sleep_time=1
run_command="go run main.go"
build_command="templ generate"
log_debug=true

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
    # macOS detected
    hash_command="shasum"
else
    # Default (linux)
    hash_command="sha1sum"
fi

log() {
  if $log_debug; then
    echo "$*"
  fi
}

log "pid $$"

# Function to gracefully kill the child processes recursively
cleanup() {
  local ppid=$1
  local children=$(ps -o pid,ppid -ax | awk -v ppid=$ppid '$2==ppid {print $1}')

  log "Child pid $children"

  for pid in $children; do
    cleanup $pid
  done;

  if [[ -n "$ppid" && -n $(ps -p "$ppid" -o args=) ]]; then

    log "Killing process $ppid"
    
    kill $ppid 2>/dev/null
    while kill -0 $ppid 2>/dev/null; do
      sleep 0.1
    done
  fi
}

stop() {
  echo ""
  echo "Cleaning up"
  if [[ -f .hotreload ]]; then
    read -r _ previous_pid < .hotreload
    rm .hotreload
  fi
  cleanup $previous_pid
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
  for file in $(find . | grep -E $file_grep | sort); do
    if [[ -z "$current_hash" ]]; then
      current_hash=$(${hash_command} "$file" | awk '{print $1}')
    else
      current_hash=$(echo -n "$current_hash$(${hash_command} "$file" | awk '{print $1}')" | ${hash_command} | awk '{print $1}')
    fi
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

    cleanup $previous_pid

    if [[ -n $build_command ]]; then
      $build_command
    fi
    # rehash after build to prevent double updates
    current_hash=$(calculate_hash)
    
    # command redirect output, err to stdout, redirect input.
    $run_command &
    command_pid=$!
    echo "$current_hash $command_pid" > .hotreload
    log "Started with pid $command_pid"
  fi

  sleep "$sleep_time"
done
