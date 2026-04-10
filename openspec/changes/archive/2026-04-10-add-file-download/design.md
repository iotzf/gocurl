## Context

gocurl is a Go-based CLI HTTP client that currently outputs response bodies directly to stdout. The tool lacks file download capability, which is a common use case when interacting with file APIs or downloading resources from URLs.

Users expect wget-style behavior when downloading files: automatic filename detection from URL, real-time progress bar, and the ability to specify custom output filenames.

## Goals / Non-Goals

**Goals:**
- Provide wget-like download experience with progress bar
- Support custom output filenames via `-O` flag
- Enable download mode via `-o` flag (rather than stdout display)
- Stream large files directly to disk to minimize memory usage
- Display real-time download statistics (speed, percentage, ETA)

**Non-Goals:**
- Multi-threaded/segmented downloading (no parallel chunks)
- Full wget compatibility (only essential download features)
- GUI progress visualization
- Background downloads

## Decisions

### 1. Progress Bar Implementation

**Decision:** Use a custom progress bar drawn using ANSI escape sequences in the terminal.

**Rationale:**
- Avoids external dependencies (keeps gocurl lightweight)
- Provides full control over display format
- Works cross-platform on modern terminals
- wget-style format: `[=====>          ] 45%  1.2MB/s  ETA 0:30`

**Alternatives considered:**
- External library (e.g., `cheggaaa/pb`): Adds dependency, less control over format
- No progress bar: Poor UX for large files

### 2. Download Mode vs Display Mode

**Decision:** Two distinct modes:
- Normal mode (`-d`): Display response to stdout (current behavior)
- Download mode (`-o`): Save to file with progress bar

**Rationale:**
- Backwards compatible with existing behavior
- Clear distinction between interactive API calls vs file downloads
- `-o` flag mirrors curl's download flag convention

### 3. Streaming Strategy

**Decision:** Use `io.Copy` with a custom `Writer` that tracks bytes written and reports progress.

**Rationale:**
- Memory efficient: constant memory usage regardless of file size
- Simple implementation using Go's standard library
- Progress callback on each write for real-time updates

### 4. Filename Handling

**Decision:** If `-O` is not provided, extract filename from URL path (URL decoding).

**Rationale:**
- wget-compatible behavior
- Fallback to `gocurl_download` if no filename in URL
- `-O` gives explicit user control

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Terminal not supporting ANSI escape sequences | Check terminal capability, fall back to simple text output |
| Content-Length header missing | Display indeterminate progress, show bytes downloaded only |
| Server doesn't support resume | Fail gracefully, restart download from beginning |
| File write permissions denied | Clear error message, exit with non-zero code |
