#!/usr/bin/env bash

function gimme() {
    local output
    output="$(command gimme "$@")"
    local exit_code="$?"
    
    if [[ "$output" == cd://* ]]; then
        cd "${output#cd://}" || return 1
    else
        [[ -n "$output" ]] && echo "$output"
        return $exit_code
    fi
}