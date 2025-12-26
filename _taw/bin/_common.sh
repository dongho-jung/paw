#!/bin/bash

# TAW Common Utilities
# Source this file: source "$(dirname "$0")/_common.sh"

# Resolve symlinks (macOS compatible)
# Usage: resolved=$(resolve_path "$path")
resolve_path() {
    local path="$1"
    while [ -L "$path" ]; do
        local dir="$(cd "$(dirname "$path")" && pwd)"
        path="$(readlink "$path")"
        [[ "$path" != /* ]] && path="$dir/$path"
    done
    echo "$(cd "$(dirname "$path")" && pwd)/$(basename "$path")"
}

# Truncate name to max 32 chars with ... in middle if needed
# Usage: display_name=$(truncate_name "$name")
truncate_name() {
    local name="$1"
    local max_len=32
    local len=${#name}

    if [ $len -le $max_len ]; then
        echo "$name"
    else
        local keep=$(( (max_len - 3) / 2 ))
        local left="${name:0:$keep}"
        local right="${name: -$keep}"
        echo "${left}...${right}"
    fi
}

# Get TAW_HOME from script location
# Assumes script is in _taw/bin/
# Usage: source _common.sh && get_taw_home
get_taw_home() {
    local script_path="$(resolve_path "${BASH_SOURCE[1]}")"
    echo "$(cd "$(dirname "$script_path")/../.." && pwd)"
}

# Debug output function
# Uses TAW_DEBUG env var
# Usage: debug "message"
debug() {
    if [ "${TAW_DEBUG:-0}" = "1" ]; then
        local script_name="$(basename "${BASH_SOURCE[1]}")"
        echo "[DEBUG $script_name] $*" >&2
    fi
}

# Tmux helper for project-specific socket
# Usage: tm "$SESSION_NAME" command args...
tm() {
    local session="$1"
    shift
    tmux -L "taw-$session" "$@"
}
