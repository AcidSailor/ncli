# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working
with code in this repository.

## Overview

`ncli` is a NETCONF command line client. It connects over SSH to a network
device (default port 830) and sends one NETCONF operation per invocation,
then prints the raw XML reply to stdout.

The whole tool is a thin [cobra](https://github.com/spf13/cobra) wrapper
around [scrapligo](https://github.com/scrapli/scrapligo)'s
`driver/netconf`. It holds no protocol logic of its own — the only local
logic is the small XML string helpers in `internal/utils`.

Twelve subcommands: `hello`, `get`, `get-config`, `edit-config`, `rpc`,
`get-schema`, `commit`, `discard-changes`, `validate`, `copy-config`,
`delete-config`, `kill-session`.

## Common commands

Tasks live in `taskfile.yml` (from the `go-scaffolds` template):

- `task run` — `go run .` (no argument passthrough; use `go run . <args>`
  directly when you need flags)
- `task test` — `go test -race ./...`
- `task lint` — **mutates files**: `golangci-lint fmt` + `run --fix`
- `task ci` — read-only gate CI enforces (`fmt --diff` + `run`)
- `task check` — `lint` then `test`
- `task build` — local snapshot via `goreleaser release --snapshot --clean`
- `task release-plan` — read-only `svu next --v0`
- `task release-tag` — cut and push the next svu tag (fires CI release)
- `task update` — pull latest go-scaffolds v1 tooling via `copier update`

Default branch is `main`. CI and release are thin callers of the template's
reusable workflows (`go-ci.yml` / `go-release.yml@v1`).

`.golangci.yaml` enables `gofumpt` (extra rules) and `golines` with
`max-len: 80`. Keep lines at or under 80 characters or `task ci` fails.
The `modernize` linter is on.

## Architecture

Three files matter; everything else is a variation on one template.

- **`main.go`** — package `main`. Injects ldflags `version`/`commit`/`date`
  into `cmd.SetVersionInfo`, calls `cmd.Execute`, prints any error to
  stderr and exits 1.
- **`cmd/ncli.go`** — the root command plus shared state. Package-level
  vars `driverOpts`, `withLock`, `withLoggingLevel`, and `logger` hold all
  global flag values. `DriverCommonOptions()` turns them into the
  `[]util.Option` slice every subcommand passes to `netconf.NewDriver`.
  `PersistentPreRunE` builds the scrapligo `logging.Instance` — but only
  when `--logging-level` is non-empty.
- **`internal/utils/utils.go`** — three pure string helpers, no
  dependencies:
  - `FlatPathToSubtreeWithValue("/a/b", "v")` → `<a><b>v</b></a>`. This is
    how a flat `--path` becomes a NETCONF subtree filter or config body.
  - `WrapWithTags(s, tag)` → `<tag>s</tag>`.
  - `NetconfStrip` trims whitespace and the trailing `]]>]]>` NETCONF 1.0
    framing delimiter.

Every subcommand file follows the same shape, and you should keep it:

```
cmd.SilenceUsage = true          // errors print bare, no usage dump
netconf.NewDriver(host, DriverCommonOptions()...)
d.Open(); defer d.Close()
if withLock != "" { d.Lock(withLock); defer d.Unlock(withLock) }
r, err := d.<Operation>(...)
if r.Failed != nil { return r.Failed }
fmt.Println(r.Result)
```

Note the two error channels: `err` is a transport or client failure,
`r.Failed` is a NETCONF `rpc-error` from the device. Both must be checked.

## Global flags

Set on the root command as persistent flags, so they apply to every
subcommand. There are **no environment variables** and no config file —
everything is a flag.

| Flag | Default | Notes |
| --- | --- | --- |
| `--host` | — | Required. Hostname or address. |
| `--port` | `830` | |
| `--username` | — | Required together with `--password`. |
| `--password` | — | Required together with `--username`. |
| `--lock` | `""` | Datastore name to lock/unlock around the call. |
| `--with-nc-version` | `1.0` | Only `1.0` or `1.1`; anything else errors. |
| `--logging-level` | `""` | `info`, `debug`, or `critical`. |

`--username`/`--password` are marked "required together", not required —
running with neither is allowed and lets SSH auth fall through.

## Conventions and gotchas

- **`DriverCommonOptions()` reads package vars.** Call it only inside
  `RunE`, after cobra has parsed flags. Calling it at init time yields
  zero values.
- **`README.md` pastes real `--help` output.** It is kept byte-identical
  to `go run . --help`. If you add or change a global flag, regenerate
  that block rather than hand-editing it — cobra realigns every column
  when the longest flag name changes.
- **Host key checking is off.** `options.WithAuthNoStrictKey()` is
  hardcoded in `DriverCommonOptions`. Intentional for lab and containerlab
  targets; there is no flag to re-enable it.
- **Transport is hardcoded to `standard`.** That is scrapligo's pure-Go
  `crypto/ssh` transport, not the `/bin/ssh` wrapper. This is what lets
  the binary build with `CGO_ENABLED=0` and ship `FROM scratch`
  (`Dockerfile.goreleaser`). Do not switch to `system` without changing
  the image.
- **`--logging-level ""` is safe.** `logger` stays `nil` and scrapligo
  substitutes a no-op logging instance in `netconf.NewDriver`. Logs go
  through `log.Print` (stderr), not `slog`.
- **`hello` is the odd one out.** It uses `d.Channel.Open()` and
  `ReadUntilPrompt` instead of `d.Open()`, so it prints the server hello
  without completing the capability exchange. It is also the only command
  that ignores `--lock`.
- **`--path /` means "everything".** An empty path split yields no
  elements, so `FlatPathToSubtreeWithValue` returns `""` — an empty
  filter, which NETCONF reads as the whole datastore.
- **`--filter-type` is only interpreted for `subtree` and `xpath`.** Any
  other value falls through the switch, leaves the filter empty, and is
  still forwarded to scrapligo unvalidated.
- **`--value` is ignored for xpath filters.** In the `xpath` branch the
  filter is `--path` verbatim.
- **`get` and `get-config` differ in how the filter reaches scrapligo.**
  `d.Get` takes the filter as its first positional argument;
  `d.GetConfig` takes it via `opoptions.WithFilter`. Not interchangeable.
- **One session per invocation.** The driver is opened and closed inside a
  single `RunE`, so nothing carries over between two `ncli` runs. That is
  why `edit-config` has `--validate`, `--commit`, and `--discard`: they
  run in that order inside the same session. `--commit` and `--discard`
  are mutually exclusive.
- **`edit-config` flag rules:** `--target` required; one of `--path` or
  `--config-file` required and mutually exclusive; `--path` requires
  `--value`.
- **No tests exist.** `task test` passes vacuously. `internal/utils` is
  pure and is the obvious place to start.
- **`main.go`, `go.mod`, and `README.md` are seed-once** in go-scaffolds.
  `task update` propagates tooling files only.
