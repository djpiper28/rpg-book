# Electron Application For RPG Book

The "backend" spawns the "frontend" (this Electron app) and hands over details of how to connect via HTTPS to the gRPC server of the backend.

## Technology

- RPG Book's frontend uses Electron, and a Vite single page application (SPA in their docs).
- `Frontend <== HTTPS / gRPC ==> Backend`
- Linting is strict and uses Eslint, and a lot of the repodog configs (made by a cool dude I used to work with - check his stuff out!)
- Pnpm is used, don't use npm

## Building, Testing, and Fun Dev Stuff

You can either use the root Makefile, or a few scripts

```sh
pnpm dev # Starts a dev server (and the backend)

pnpm test

pnpm lint
pnpm lint --fix # Tries to fix the errors

pnpm storybook
pnpm build-storybook

pnpm build # Only builds the electron app image
```

## Updating Electron???

You need to change:

- ./nix-bins.sh
- ./package.json (deps and scripts)
