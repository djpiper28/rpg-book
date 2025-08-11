# What is used

## Also read the README.md as I am not going to type things twice.

This app uses Github Actions, Vite, Vitest Go, Electron, Electron Builder, Typescript, React, React Router, Makefile, Sqlite3, gRPC, and Mantine as the key technologies.

- Always use the makefile to execute commands (see makefile and readme)

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

# RPG Book take code and architecture quality seriously

1. Comments that point out the obvious are banned

```js
// Example BANNED comment
console.log("something broke", "error", err);
```

2. Installing random npm/go packages that do simple things is banned

3. Random `console.logs` / `printf` statements are banned

- Use the charm logger in Go i.e:

```go
log.Info("Something has happened", "what", "a man fell into the river in lego city", "when", time.Now())
```

- Use getLogger() and the same conventions in the Electron App

4. Do not touch things that you do not need to

5. Run `make format -j` when you are done with things

6. Run `make test -j` whenever you think something may be broken

7. Do not ramble in commit messages. Use conventional commits and NEVER PUSH. (also never commit and try to force users to do actual code review)

8. Do basic security. SQL injection, never trust what the user inputs, no exec, etc...

9. Read config and setup files before doing stuff.

# Installing new packages

Check for malware!!! There are a lot of malicious packages IN ALL ECO-SYSTEMS at the moment. They use similar names to useful packages i.e: `core-js-patched`. They may crypto mine, delete things, be ransomware, steal cookies, rootkits, etc...

The best way to avoid these attacks is to NOT INSTALL THINGS unless required. And if needed check the pages downloads, the linked Github repo, and use common sense. If it looks sus DO NOT USE IT.
