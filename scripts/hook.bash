#!/bin/bash

cmd="$(dirname "${BASH_SOURCE[0]}")/sdkm"

sdkm-export-env() {
	$cmd env go 2>&1 | while IFS='' read -r line; do
		export "${line}" 1>/dev/null 2>&1
	done
}

trap sdkm-export-env INT
#source <($cmd completion bash)
