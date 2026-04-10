## ADDED Requirements

### Requirement: File download mode
The system SHALL provide a download mode that saves the response body to a file instead of displaying it to stdout.

#### Scenario: Download with explicit output filename
- **WHEN** user provides `-O <filename>` flag with a URL
- **THEN** the system SHALL download the file and save it to the specified filename
- **AND** the system SHALL display a wget-style progress bar during download

#### Scenario: Download with auto-detected filename
- **WHEN** user provides `-o` flag with a URL but no explicit filename
- **THEN** the system SHALL extract the filename from the URL path
- **AND** if no filename is found in the URL, use `gocurl_download` as the default filename
- **AND** the system SHALL display a wget-style progress bar during download

#### Scenario: Download with custom content type
- **WHEN** user provides `-d <data>` along with `-o` flag
- **THEN** the system SHALL send the POST request and download the response body

### Requirement: Progress bar display
The system SHALL display a real-time progress bar during file downloads, showing download percentage, speed, and ETA.

#### Scenario: Progress bar shows completion percentage
- **WHEN** a download is in progress with known Content-Length
- **THEN** the progress bar SHALL display the percentage of bytes downloaded
- **AND** the progress bar SHALL use a visual block representation (e.g., `[====>    ]`)

#### Scenario: Progress bar shows download speed
- **WHEN** a download is in progress
- **THEN** the progress bar SHALL display the current download speed
- **AND** speed SHALL be formatted appropriately (B/s, KB/s, MB/s, GB/s)

#### Scenario: Progress bar shows ETA or elapsed time
- **WHEN** a download is in progress
- **THEN** if total size is known, display estimated time remaining
- **AND** if total size is unknown, display elapsed time since download started

#### Scenario: Progress bar shows downloaded bytes
- **WHEN** a download is in progress
- **THEN** the progress bar SHALL display the number of bytes downloaded
- **AND** if total size is known, display as `downloaded / total` format

### Requirement: Large file streaming
The system SHALL stream downloaded files directly to disk to minimize memory usage.

#### Scenario: Large file download does not exhaust memory
- **WHEN** user downloads a file larger than available RAM
- **THEN** the system SHALL write chunks to disk as they are received
- **AND** memory usage SHALL remain constant regardless of file size

### Requirement: Resume partial downloads
The system SHALL support resuming downloads when the server supports the Content-Range header.

#### Scenario: Resume supported by server
- **WHEN** user runs the same download command for a partially downloaded file
- **AND** the server supports the Content-Range header
- **THEN** the system SHALL resume downloading from where it left off
- **AND** the progress bar SHALL reflect the resumed position

#### Scenario: Resume not supported by server
- **WHEN** user runs the same download command for a partially downloaded file
- **AND** the server does NOT support the Content-Range header
- **THEN** the system SHALL restart the download from the beginning
- **AND** overwrite the partial file

### Requirement: File overwrite protection
The system SHALL prompt or handle existing files appropriately.

#### Scenario: Output file already exists
- **WHEN** the output filename already exists
- **THEN** the system SHALL overwrite the file without prompting (matching wget behavior)
