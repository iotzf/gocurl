## 1. CLI Changes

- [x] 1.1 Add `-o, --download` flag to enable download mode
- [x] 1.2 Add `-O, --output` flag to specify output filename
- [x] 1.3 Implement filename extraction from URL path as fallback
- [x] 1.4 Pass download mode and output filename to httpclient

## 2. Progress Bar Implementation

- [x] 2.1 Create `ProgressBar` struct with ANSI escape sequence rendering
- [x] 2.2 Implement `Write` method to update progress on each chunk
- [x] 2.3 Add speed calculation and formatting (B/s, KB/s, MB/s, GB/s)
- [x] 2.4 Add ETA/elapsed time calculation
- [x] 2.5 Add percentage and visual progress block rendering
- [x] 2.6 Support terminals that don't support ANSI (fallback to simple text)

## 3. HTTP Client Download Support

- [x] 3.1 Add `DownloadMode` and `OutputFilename` fields to RequestOptions
- [x] 3.2 Create streaming file writer with progress callback
- [x] 3.3 Implement `io.Copy` with progress tracking for download mode
- [x] 3.4 Add Content-Range header support for resume capability
- [x] 3.5 Handle partial file exists for resume scenario

## 4. Integration and Testing

- [x] 4.1 Test download with explicit filename
- [x] 4.2 Test download with auto-detected filename from URL
- [x] 4.3 Test progress bar displays correctly during download
- [x] 4.4 Test large file download (verify memory efficiency)
- [x] 4.5 Test download with POST data
