#!/bin/bash

cd ./node_modules/electron/dist/ || exit 1
rm electron
ln -s "$(whereis electron | cut -f 2 -d ' ')"
