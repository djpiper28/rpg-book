#!/bin/bash

NODE_MODULES="$(pwd)/node_modules"

cd "$NODE_MODULES/electron/dist/" || exit 1
rm electron
ln -s "$(whereis electron | cut -f 2 -d ' ')"

cd "$NODE_MODULES/.bin" || exit 1
rm protoc-gen-js
ln -s "$(whereis protoc-gen-js | cut -f 2 -d ' ')"
