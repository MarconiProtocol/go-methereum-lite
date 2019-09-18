#!/bin/sh

set -e

if [ ! -f "build/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
root="$PWD"
ethdir="$workspace/src/github.com/MarconiProtocol"
if [ ! -L "$ethdir/go-methereum-lite" ]; then
    mkdir -p "$ethdir"
    cd "$ethdir"
    ln -s ../../../../../. go-methereum-lite
    cd "$root"
fi

# Set up the environment to use the workspace.
GOPATH="$workspace"
export GOPATH

# Run the command inside the workspace.
cd "$ethdir/go-methereum-lite"
PWD="$ethdir/go-methereum-lite"

# Launch the arguments with the configured environment.
exec "$@"
