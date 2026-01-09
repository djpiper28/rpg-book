# RPG Book

## Getting Started

Go to the [builds page](https://github.com/djpiper28/rpg-book/actions/workflows/release.yml?query=branch%3Amain), select a build and download the zip file for your OS.
On Windows and Mac this is an installer that you use; on Linux this is an app image. Once installed, if necessary, you can then run the program and it should just work.

---

## Dev Stuff

This is a work in progress, stayed tuned for updates.

> An all in one tool for tabletop RPG DM-ing.

[![test](https://github.com/djpiper28/rpg-book/actions/workflows/test.yml/badge.svg)](https://github.com/djpiper28/rpg-book/actions/workflows/test.yml)
[![build](https://github.com/djpiper28/rpg-book/actions/workflows/build.yml/badge.svg)](https://github.com/djpiper28/rpg-book/actions/workflows/build.yml)
[![release](https://github.com/djpiper28/rpg-book/actions/workflows/release.yml/badge.svg)](https://github.com/djpiper28/rpg-book/actions/workflows/release.yml)

### Architecture

A backend talks to a frontend, for the desktop app that is a Go server talking to an Electron app via gRPC. Data is stored in Sqlite3 files, each project gets its own
database, and global settings their own database as well.

### `common/` contains Go code shared by the web app and the desktop app.

#### `common/search/`

This uses [Go Peg](https://github.com/pointlander/peg) to generate the parser.

- https://github.com/pointlander/peg/blob/main/README.md
- https://github.com/pointlander/peg/blob/main/peg.peg
- https://github.com/pointlander/peg/blob/main/docs/peg-file-syntax.md
- https://github.com/pointlander/peg/blob/main/docs/links.md

### `desktop_client/` contains Typescript code for the desktop app

## Building, Testing, and Dev Stuff

```sh
nix-shell shell.nix # You may with to use the `--pure` flag if you have issues
make -j # Builds the backend and frontend
# The frontend is in ./desktop_client/app/release/
# The launcher is in ./desktop_client/launcher/launcher
# See the release target for how to use this.

make test -j # Runs the codegen, and tests
make format -j # Runs all formatter scripts
make fuzz -j # Runs the fuzzers sequentially, this takes a while

make dev -j # Starts a hot-reloading dev version of the app
make release -j # Starts a release build of the app
```
