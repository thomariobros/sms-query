#!/bin/bash
set -e

if [[ "$(gofmt -s -l . | wc -l)" -gt 0 ]]; then
  exit 1
fi
