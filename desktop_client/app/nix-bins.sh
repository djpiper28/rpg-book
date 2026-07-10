#!/bin/bash

pushd ./node_modules/.pnpm/electron@*/node_modules/@electron/get/dist/esm
rm -f electron
ln -s "$(whereis electron | cut -f 2 -d ' ')"
popd
