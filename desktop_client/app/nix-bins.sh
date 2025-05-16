#!/bin/bash

NODE_MODULES="$(pwd)/node_modules"

cd "$NODE_MODULES/electron/dist/" || exit 1
rm electron
ln -s "$(whereis electron | cut -f 2 -d ' ')"
