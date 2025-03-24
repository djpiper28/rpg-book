# RPG Book

## Architecture

There is a Go monolith that uses SQLite3 when ran on the desktop app, and Postgres when used in the web app.

The web app will display different addons that can be installed, and handle cloud backups.

## Folder Structure

### `common/` contains Go code shared by the web app and the desktop app.

### `desktop_client/` contains Typescript code for the desktop app

### `webapp/frontend` contains Typescript code for the web app

### `webapp/backend` contains Go code for the web app
