# rhun extension template (Go)

A rhun extension is a program that rhun starts and talks to over its standard input and output,
one text line per message. This repository is the smallest useful one: commands with a toast and
an edit, a key, a menu item and a status item, in Go, with the release workflow that builds it for Linux, macOS and Windows. Rust is the
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
   `tests/basic.txt`, and change what `Invoke` does. `Invoke` receives each run with the
   command's name (`inv.Name()` is `say`; `inv.Command` the full `hello.say`), where it came from
   and what that place passes (below). There is no request
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

## Keys, menus, the status bar and language keys

`extension.ini` declares more than commands (protocol 1.1, rhun 0.18):

- `[keys]`: `ctrl+k ctrl+h = greet world when editorFocus` binds a chord to `greet` with the
  argument `world`, only while text has the focus. One stroke or two; a bare letter needs a
  `when` that keeps it away from typing.
- `[menu.editor]` (also `explorer`, `tab`, `terminal`, `editor_title` and the views' menus):
  `greet.when = resourceExt in '.txt .md'` shows **Greet** in the editor's context menu for those
  files; `greet.args` is what the item passes.
- `[status]`: `hello.text = $(play) Hello` puts an item in the status bar, whose click runs
  `hello.command` with `hello.args`.
- `[language.Shell]`: language keys for a grammar rhun has (`continue = #` continues a comment
  on Enter; also brackets, auto-closing pairs, indentation rules, comments, folding markers).

`main.go`'s `Invoke` sees each run's `inv.From` (`palette`, `key`, `menu`, `ref`) and `inv.Body`:
the binding's arguments, a menu's target line and the item's arguments (`inv.Target()`), a
command reference's arguments. It needs the `command-args` feature, which `main` asks for.
The Extensions page lists all of these under **Contributions**, with any line rhun dropped.
`tests/basic.txt` runs `greet` from the palette, its chord, the menu and the status item.

The SDK's protocol 1.1 additions come with `rhunext` v1.4.0. To try an SDK change before it is
tagged, `RHUN_EXT_SDK=/path/to/extension-protocol sh tests/run.sh` (and `tests/check.sh`) build
against that checkout through a `go.work` in their temporary folder.

## One module for every system (WebAssembly)

Instead of three binaries, an extension can ship one WebAssembly module that rhun runs itself on
Linux, macOS and Windows, held to its `[needs]` on every system. In `extension.ini` replace the
three `run.TARGET` lines with

```ini
run.wasm32-wasip1 = bin/hello.wasm
```

`GOOS=wasip1 GOARCH=wasm go build -o bin/hello.wasm .` builds it (`rhunext` reads standard input
non-blocking there, so goroutines keep running), `tests/run.sh` tests it when the manifest names
no program for your machine, and the release workflow packs it into every archive. Keeping the
`run.TARGET` lines as well gives those systems the native binary (faster at heavy work) and the
module to any other. What the module may use (WASI preview 1, no sockets, no processes): section
"WebAssembly extensions" of [protocol 1.1](https://rhun.app/extensions/v1.1/).

## The scripts variant

`scripts/main.sh` and `scripts/main.ps1` are the same extension as a shell script and its
PowerShell twin, for extensions that only glue a tool to the editor. Their manifest uses
`runtime = sh` with `run = scripts/main.sh` (Linux and macOS; `runtime = powershell` on Windows)
instead of the `run.*` binaries, and ships one archive for all targets.

## Versions, development installs and the index

- Each tag installs on its own: `OWNER/NAME@0.1.0` in **Install from GitHub** (or the link
  `rhun://extension/OWNER/NAME@0.1.0`) installs that release; the name alone installs the latest
  release's exact tag. rhun keeps the version an update replaces, so users can **Roll back**, and
  a new version that fails its first start (no `hello` within 5 seconds) rolls back by itself.
- While you work, the Extensions page's **Install** menu has **Link a Folder (development)**: rhun
  runs the extension from your folder (build `bin/` first), so rebuild and **Restart** is the loop.
  **From a Folder** and **From an Archive** install a copy or a release `.tar.gz` you built.
- To be listed under **Browse**, and updated by itself on users' machines, open a pull request to
  [rhun-editor/extensions](https://github.com/rhun-editor/extensions) adding `OWNER/NAME.ini` that
  pins a release's tag and digests; its checks verify the digests and this workflow's attestation.
  An update that asks for more access (`[needs]`) always asks users first.
- A pack installs several extensions together: an extension without a program whose
  `extension.ini` says `pack = OWNER/NAME …` (at most 16). Release it like any other, with
  `NAME-TARGET.tar.gz` (its `extension.ini`, `README.md` and `LICENSE`) for each target and
  `sha256sums.txt`.

## State and environment

Your extension's folder is read-only once installed. Write to `$RHUN_STATE`, which is also its
`HOME` and holds `TMPDIR` and the common caches. `$RHUN_PROJECT` is the project folder when there
is one. Standard error goes to a log file rhun keeps under 1 MB.
