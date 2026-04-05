## Context

nestbird runs `netbird up` to reconnect when disconnected. Some environments require a NetBird setup key for authentication on `netbird up` — without it, the command fails silently or prompts interactively. The setup key is sensitive and should be stored in a file (not passed as a plain string), which NetBird supports natively via `--setup-key-file`.

Currently, `NetBird` is a stateless struct and `NewNetBird()` takes no arguments. `main.go` has no flag parsing.

## Goals / Non-Goals

**Goals:**
- Pass `--setup-key-file=<path>` to `netbird up` when a path is configured
- Accept the path via CLI flag `--setup-key-file` or env var `NETBIRD_SETUP_KEY_FILE`, with CLI flag taking precedence
- Keep the `NetBirdClient` interface unchanged (callers are unaffected)

**Non-Goals:**
- Reading or validating the contents of the setup key file
- Supporting inline setup keys (only file-based)
- Any other new `netbird up` flags

## Decisions

### Bake path into `NetBird` at construction, not per-call

The path is static for the lifetime of the process. Rather than threading it through the `NetBirdClient.Up(ctx)` signature (which would require an interface change and affect all mocks/tests), it is stored on the struct at `NewNetBird(setupKeyFile string)` construction time.

**Alternatives considered:**
- Passing path as `Up(ctx, setupKeyFile string)` — rejected because it changes the interface and forces all callers/mocks to update with no benefit (the value never changes at runtime).
- Config struct — overkill for a single optional string.

### CLI flag wins over env var

Standard convention (`flag` package is parsed first; if flag is non-empty, it is used; otherwise fall back to `os.Getenv`).

### No validation of path existence at startup

The file's existence is only meaningful at reconnect time. Validating at startup would add complexity and could cause the watcher to refuse to start on a transient filesystem issue. `netbird up` will surface the error naturally.

## Risks / Trade-offs

- **Secret exposure in process argv** → Mitigated: we pass `--setup-key-file=<path>`, not the key value itself. The path appears in `ps` output but the secret does not.
- **File unreadable at reconnect time** → `netbird up` returns a non-zero exit code; watcher handles this via existing backoff/retry logic. No special handling needed.
- **Breaking `NewNetBird()` call signature** → `main.go` is the only caller; update is trivial and contained.
