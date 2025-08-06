# RPG Book

This is a work in progress, stayed tuned for updates.

> An all in one tool for tabletop RPG DM-ing.

[![test](https://github.com/djpiper28/rpg-book/actions/workflows/test.yml/badge.svg)](https://github.com/djpiper28/rpg-book/actions/workflows/test.yml)
[![build](https://github.com/djpiper28/rpg-book/actions/workflows/build.yml/badge.svg)](https://github.com/djpiper28/rpg-book/actions/workflows/build.yml)
[![release](https://github.com/djpiper28/rpg-book/actions/workflows/release.yml/badge.svg)](https://github.com/djpiper28/rpg-book/actions/workflows/release.yml)

## Architecture

A backend talks to a frontend, for the desktop app that is a Go server talking to an Electron app via gRPC. Data is stored in Sqlite3 files, each project gets its own database, and global settings their own database as well.

## Folder Structure

### `common/` contains Go code shared by the web app and the desktop app.

### `desktop_client/` contains Typescript code for the desktop app

## Building, Testing, and Dev Stuff

```sh
make -j # Builds the backend and frontend
# The frontend is in ./desktop_client/app/release/
# The launcher is in ./desktop_client/launcher/launcher
# See the release target for how to use this.
# TODO: create a production ready script for starting the app

make test -j # Runs the codegen, and tests
make format -j # Runs all formatter scripts

make dev -j # Starts a hot-reloading dev version of the app
make release -j # Starts a release build of the app
```
