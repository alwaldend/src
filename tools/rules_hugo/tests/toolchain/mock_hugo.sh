#!/bin/sh
if [ "${1:-}" = "env" ]; then
    echo "HUGO=mock-hugo"
    exit 0
fi
echo "mock hugo build"
