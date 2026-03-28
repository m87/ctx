#!/bin/sh
set -eu

ctx serve --addr :8080 &

exec nginx -g 'daemon off;'
