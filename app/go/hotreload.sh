#!/bin/bash

### Parameters
source_file_grep=".+(go|rs|templ)$"       # example '.+(go|rs|templ)$'
ws_refresh_file_grep=".+(css|html|js)$"   # example '.+(css|html|js)$'
run_command="go run main.go"              # example 'go run main.go'
build_command=""                          # example 'templ generate'
sleep_time=1
log_debug=true
enable_ws_refresh=true
ws_port=7878
###

if $enable_ws_refresh && [[ -z $(which websocat) ]]; then
  echo "Needs websocat to enable browser ws refresh"
  echo "Install with aptly or brew"
  echo "Or set enable_ws_refresh to false"
  exit 0
fi

echo "Hotreload settings:"
echo "source file grep:     $source_file_grep"
echo "ws refresh file grep: $ws_refresh_file_grep"
echo "run command:          $run_command"
echo "build command:        $build_command"
echo "sleep time (sec):     $sleep_time"
echo "hotreload pid:        $$"
echo "log debug:            $log_debug"
echo "enable ws refresh:    $enable_ws_refresh"
echo "ws port:              $ws_port"
echo
echo

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


# Function to calculate cumulative hash of specified files
calculate_hash() {
  local file_grep=$1
  local current_hash="x"
  for file in $(find . | grep -E $file_grep | sort); do
    if [[ -f $file ]]; then
      current_hash=$(echo -n "$current_hash$(${hash_command} "$file" | awk '{print $1}')" | ${hash_command} | awk '{print $1}')
    fi
  done
  echo "$current_hash"
}

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

# Cleanup function for child processes and websocat
stop() {
  echo ""
  echo "Cleaning up"
  if [[ -f .hotreload ]]; then
    read -r _ _ previous_pid < .hotreload
    rm .hotreload
  fi
  cleanup $previous_pid
  if [[ -n $ws_pid ]]; then
    kill $ws_pid
  fi
  exit 0
}

# Trap TERM signal and call cleanup function
trap stop TERM INT


# Startup
echo Startup
# Starting websocat for hot-reloading in browser
if $enable_ws_refresh; then
  if $log_debug; then
    websocat -t ws-l:127.0.0.1:$ws_port broadcast:mirror: &
  else
    websocat -t ws-l:127.0.0.1:$ws_port broadcast:mirror: 2&>1 1>/dev/null &
  fi
  # store websocat pid for stopping
  ws_pid=$!
  log "websocket opened at ws://127.0.0.1:$ws_port"
fi
current_ws_refresh_source_hash=$(calculate_hash $ws_refresh_file_grep)
if [[ -n $build_command ]]; then
  $build_command
fi
current_source_hash=$(calculate_hash $source_file_grep)
$run_command &
command_pid=$!
echo "$current_source_hash $current_ws_refresh_source_hash $command_pid" > .hotreload
echo End Startup
# End Startup

# Main loop
echo Entering Main Loop
while true; do
  current_source_hash=$(calculate_hash $source_file_grep)
  current_ws_refresh_source_hash=$(calculate_hash $ws_refresh_file_grep)

  if [[ -f .hotreload ]]; then
    read -r source_hash ws_refresh_hash previous_pid < .hotreload
  fi

  # Source files hot reloading
  if [[ "$current_source_hash" != "$source_hash" ]]; then
    echo "Hot-Reloading..."

    cleanup $previous_pid

    if [[ -n $build_command ]]; then
      $build_command
    fi
    # rehash after build to prevent double updates
    current_source_hash=$(calculate_hash $source_file_grep)
    
    # command redirect output, err to stdout, redirect input.
    $run_command &
    command_pid=$!
    echo "$current_source_hash $current_ws_refresh_source_hash $command_pid" > .hotreload
    log "Run command started with pid $command_pid, sleeping 1s for startup"
    sleep 1
    if $enable_ws_refresh; then
      echo refresh | websocat -1 ws://127.0.0.1:$ws_port
    fi
  fi

  # Browser refresh
  if $enable_ws_refresh && [[ "$current_ws_refresh_source_hash" != "$ws_refresh_hash" ]]; then
    echo refresh | websocat -1 ws://127.0.0.1:$ws_port
    log "Refresh command send to websocket clients"
    echo "$current_source_hash $current_ws_refresh_source_hash $command_pid" > .hotreload
  fi

  sleep "$sleep_time"
done
