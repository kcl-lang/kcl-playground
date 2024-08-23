#!/bin/bash

# Start pluto locally
TMP_OUTPUT_FILE=/tmp/pluto_output.log
pluto run >$TMP_OUTPUT_FILE 2>&1 &
PID1=$!

# Wait for Pluto to start the project completely
while :; do
    if grep -q "Successfully applied!" $TMP_OUTPUT_FILE; then
        echo "The project has been successfully started."
        break
    else
        echo "Waiting for Pluto to start the project..."
        sleep 1
    fi
done

# Get the project name from package.json
PROJECT_NAME=$(grep '"name":' package.json | awk -F '"' '{print $4}')
echo "Project name: $PROJECT_NAME"

# Set environment variables
PORT=$(lsof -Pan -p $PID1 -i | grep LISTEN | awk '{print $9}' | sed 's/.*://')
export PLUTO_PROJECT_NAME=$PROJECT_NAME
export PLUTO_STACK_NAME=local_run
export PLUTO_PLATFORM_TYPE=Simulator
export PLUTO_SIMULATOR_URL=http://localhost:$PORT

# Run tests
print_separator() {
    local message=$1
    local message_len=${#message}

    local width=$(tput cols)
    local separator=$(printf '=%.0s' $(seq 1 $(((width - message_len - 4) / 2))))
    local bold=$(tput bold)

    printf "\033[34m${bold}${separator}= %s =${separator}\033[0m\n" "$message"
}

# Output Pluto logs, which might contain useful information
tail -f $TMP_OUTPUT_FILE -n 0 &
PID2=$!

# Execute tests in the app directory
print_separator "Executing test files in the app directory"
python3 -m pytest -s -q --no-header app

# Execute tests within the app/main.py file
print_separator "Executing tests within the app/main.py file"
python3 -m pytest -s -q --no-header app/main.py

# Cleanup
kill $PID1
wait $PID1
kill $PID2
