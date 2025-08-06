# RPG Book

> An all in one tool for tabletop RPG DM-ing.

[![test](https://github.com/djpiper28/rpg-book/actions/workflows/test.yml/badge.svg)](https://github.com/djpiper28/rpg-book/actions/workflows/test.yml)
[![build](https://github.com/djpiper28/rpg-book/actions/workflows/build.yml/badge.svg)](https://github.com/djpiper28/rpg-book/actions/workflows/build.yml)

## Architecture

A backend talks to a frontend, for the desktop app that is a Go server talking to an Electron app via gRPC. Data is stored in Sqlite3 files, each project gets its own database, and global settings their own database as well.

## Folder Structure

### `common/` contains Go code shared by the web app and the desktop app.

### `desktop_client/` contains Typescript code for the desktop app

## Building, Testing, and Dev Stuff

```sh
make build -j # Codegen, frontend, and backend builds
make test -j # Runs the codegen, and tests
make format -j # Runs all formatter scripts
```
