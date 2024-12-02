#!/bin/bash

# Change to the correct directory and run the make command
cd server && make watch

# You can also run other services or processes after this
exec "$@"  # This allows other CMD commands to be run (if needed)