## Why

gocurl currently only displays response bodies in the terminal. Users need to download files directly, similar to wget, with real-time progress visualization showing download speed, percentage, and time remaining.

## What Changes

- Add `-O, --output` flag to specify output filename
- Add `-o, --download` flag to enable download mode (rather than display mode)
- Implement wget-style progress bar showing:
  - Download percentage with progress bar visualization
  - Download speed (bytes/sec, KB/s, MB/s)
  - Time remaining or elapsed time
- Stream file to disk instead of memory for large files
- Support resuming partial downloads via `Content-Range` header (if server supports)

## Capabilities

### New Capabilities

- `file-download`: Download files from URLs with wget-style progress bar, supporting output filename specification, large file streaming, and resume capability

## Impact

- **CLI**: New flags in `main.go` (`-O`, `--output` and `-o`, `--download`)
- **HTTP Client**: Modified `httpclient.go` to support streaming downloads with progress callbacks
- **Dependencies**: None (using standard library for progress bar rendering)
