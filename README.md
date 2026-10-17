# rhun extension template (Go)

A rhun extension is a program that rhun starts and talks to over its standard input and output,
one text line per message. This repository is the smallest useful one: two commands, a toast and
an edit, in Go, with the release workflow that builds it for Linux, macOS and Windows. Rust is the
first choice for new extensions: see [template-rust](https://github.com/rhun-editor/template-rust).

Written for people and for coding agents: [AGENTS.md](AGENTS.md) is the working guide, and
`.claude/skills/rhun-extension` is a skill Claude Code picks up in this folder (Codex: install
the `rhun` plugin from [rhun-editor/skills](https://github.com/rhun-editor/skills)).

The protocol: https://rhun.app/extensions/v1/ (the `rhunext` Go package does the framing).

## Make your own

1. Copy this repository under your GitHub account as `OWNER/NAME`.
2. In `extension.ini` set `name = NAME`, `owner = OWNER`, `title`, your commands and what the
   extension needs (`[needs]`). The `name` and `owner` must equal the repository's owner and
   name: rhun installs into `~/.config/rhun/extensions/OWNER/NAME/`.
3. Rename `hello` in `main.go`, `extension.ini` and the comment at the top of
   `tests/basic.txt`, and change what `Run` does. `Run` receives the command's full
   name, `NAME.key` (`hello.say`, not `say`), as the `run` event carries it. There is no request
   for the cursor position: an extension that acts on the word under the cursor subscribes to
   `cursor` in `[subscribe]` and remembers the last position per document.
4. Test headlessly: `tests/run.sh` builds the extension, lays it out the way rhun installs it,
   drives rhun with `tests/basic.txt` and compares the output with `tests/basic.expected`
   (section 10 of the spec lists the script commands and what they print; the state folder in
   such a run is `$XDG_STATE_HOME/rhun/extensions/OWNER--NAME`, so a test can seed files there).
5. Tag `vX.Y.Z` (the same as `version` in the manifest). The release workflow builds the three
   targets, attests them and publishes `NAME-TARGET.tar.gz` and `sha256sums.txt`, which
   **Install from GitHub** takes. `tests/check.sh` runs the conformance checker.
6. Open a pull request to `rhun-editor/extensions` to be listed in the catalog.

## The scripts variant

`scripts/main.sh` and `scripts/main.ps1` are the same extension as a shell script and its
PowerShell twin, for extensions that only glue a tool to the editor. Their manifest uses
`runtime = sh` with `run = scripts/main.sh` (Linux and macOS; `runtime = powershell` on Windows)
instead of the `run.*` binaries, and ships one archive for all targets.

## State and environment

Your extension's folder is read-only once installed. Write to `$RHUN_STATE`, which is also its
`HOME` and holds `TMPDIR` and the common caches. `$RHUN_PROJECT` is the project folder when there
is one. Standard error goes to a log file rhun keeps under 1 MB.
