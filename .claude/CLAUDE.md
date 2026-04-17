# NLM CLI - NotebookLM Command Line Interface

## Project Overview

This is a Go CLI for interacting with Google NotebookLM. It uses Google's internal batchexecute RPC protocol to communicate with NotebookLM's backend.

Based on [tmc/nlm](https://github.com/tmc/nlm) upstream with fork-specific patches for Drive URL support, macOS auth fixes, and `--format plain` output.

## Architecture

```
cmd/nlm/main.go              - CLI entry point, command routing
internal/notebooklm/api/     - High-level API client
internal/notebooklm/rpc/     - Low-level RPC client
internal/batchexecute/       - Google batchexecute protocol
internal/beprotojson/        - Proto <-> batchexecute JSON marshaling
internal/nlmmcp/             - MCP server implementation
internal/auth/               - Browser cookie extraction
gen/
  gen/method/                - RPC argument encoders
  gen/service/               - Service clients
  gen/notebooklm/v1alpha1/   - Protocol buffer definitions
```

## Build & Test

Requires libopus/opusfile/libogg for the interactive audio module:

```bash
# Build (with nix dependencies)
nix-shell -p libopus opusfile libogg pkg-config --run "go build ./cmd/nlm"

# Run tests
go test ./...

# Install locally
go install ./cmd/nlm
```

## Authentication

NLM uses browser cookies for authentication. Run `nlm auth` to set up.

Chrome always runs non-headless during auth to enable macOS keychain access for cookie decryption.

## Fork-Specific Features

- **Google Drive URLs**: `nlm add <id> https://drive.google.com/file/d/FILE_ID/view` works with Drive-specific payload format
- **PDF / binary file uploads**: `nlm add <id> file.pdf` uses Google Resumable Upload Protocol, HAR-verified flow (register via `o4cbdc` → upload init → upload bytes)
- **--format plain**: `nlm --format plain generate-chat <id> "question"` produces clean text output for piping
- **Release workflow**: `.github/workflows/release.yml` builds multi-platform binaries on tag push

## PDF Upload Flow (HAR-verified)

The resumable upload has a specific order that must be preserved (see `internal/notebooklm/api/client.go`):

1. `registerFileSource` (RPC `o4cbdc`) with `[[[filename]], projectID, [2], [1, nil*9, [1]]]` — server creates tentative queue entry and returns the source_id
2. `startResumableUpload` — POST to `https://notebooklm.google.com/upload/_/?authuser=N` with **plain JSON** body `{"PROJECT_ID","SOURCE_ID","SOURCE_NAME"}` (NOT base64-encoded). Returns `X-Goog-Upload-Url`
3. `uploadFileBytes` — POST raw binary bytes with `X-Goog-Upload-Command: upload, finalize`
4. `processFileSource` — optional trigger for document guide generation (non-fatal if unavailable)
