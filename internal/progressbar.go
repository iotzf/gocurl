package httpclient

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// ProgressBar provides wget-style progress display
type ProgressBar struct {
	total      int64
	downloaded int64
	startTime  time.Time
	terminal   bool
	filename   string
}

// NewProgressBar creates a new ProgressBar
func NewProgressBar(total int64, filename string) *ProgressBar {
	// Detect if we're writing to a terminal
	_, err := os.Stdout.Stat()
	terminal := err == nil

	pb := &ProgressBar{
		total:      total,
		downloaded: 0,
		startTime:  time.Now(),
		terminal:   terminal,
		filename:   filename,
	}

	if pb.terminal {
		fmt.Fprintf(os.Stdout, "Saving to: '%s'\n", filename)
	}
	return pb
}

// Write implements io.Writer, tracking progress
func (pb *ProgressBar) Write(p []byte) (n int, err error) {
	n = len(p)
	pb.downloaded += int64(n)
	return n, nil
}

// SetDownloaded sets the current downloaded bytes (for resume support)
func (pb *ProgressBar) SetDownloaded(d int64) {
	pb.downloaded = d
}

// Finish prints final progress and newline
func (pb *ProgressBar) Finish() {
	pb.printProgress()
	if pb.terminal {
		fmt.Fprintln(os.Stdout)
	}
	os.Stdout.Sync()
}

func (pb *ProgressBar) printProgress() {
	downloaded := pb.downloaded
	if downloaded < 0 {
		downloaded = 0
	}

	elapsed := time.Since(pb.startTime)
	var speed float64
	if elapsed.Seconds() > 0 && downloaded > 0 {
		speed = float64(downloaded) / elapsed.Seconds()
	}

	// Format speed
	speedStr := formatSpeed(speed)

	// Calculate percentage and ETA
	var etaStr string
	var percent float64
	barLen := 40
	var barStr string

	if pb.total > 0 {
		percent = float64(downloaded) / float64(pb.total) * 100
		if percent > 100 {
			percent = 100
		}
		filled := int(float64(barLen) * percent / 100)
		barStr = "[" + strings.Repeat("=", filled)
		if filled < barLen {
			barStr += ">"
			if filled < barLen-1 {
				barStr += strings.Repeat(" ", barLen-filled-1)
			}
		}
		barStr += "]"

		remaining := float64(pb.total - downloaded)
		if speed > 0 && remaining > 0 {
			eta := time.Duration(remaining / speed * float64(time.Second))
			etaStr = formatDuration(eta)
		} else {
			etaStr = "--:--"
		}
	} else {
		// Unknown total size - show indeterminate progress
		barStr = "[" + strings.Repeat(" ", barLen) + "]"
		etaStr = formatDuration(elapsed)
	}

	// Format sizes
	downloadedStr := formatSize(downloaded)
	var totalStr string
	if pb.total > 0 {
		totalStr = formatSize(pb.total)
	} else {
		totalStr = "unknown"
	}

	// Build progress string
	progress := fmt.Sprintf("\r%s %5.1f%% %7s/s %8s  %s/%s   ",
		barStr, percent, speedStr, etaStr, downloadedStr, totalStr)

	// Clear line and print
	fmt.Fprint(os.Stdout, progress)
}

func formatSize(bytes int64) string {
	if bytes < 0 {
		bytes = 0
	}
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%dB", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func formatSpeed(bytesPerSec float64) string {
	if bytesPerSec < 0 {
		bytesPerSec = 0
	}
	if bytesPerSec < 1024 {
		return fmt.Sprintf("%.0fB", bytesPerSec)
	}
	if bytesPerSec < 1024*1024 {
		return fmt.Sprintf("%.1fK", bytesPerSec/1024)
	}
	if bytesPerSec < 1024*1024*1024 {
		return fmt.Sprintf("%.1fM", bytesPerSec/1024/1024)
	}
	return fmt.Sprintf("%.1fG", bytesPerSec/1024/1024/1024)
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	if d.Hours() >= 1 {
		return fmt.Sprintf("%d:%02d:%02d", int(d.Hours()), int(d.Minutes())%60, int(d.Seconds())%60)
	}
	if d.Minutes() >= 1 {
		return fmt.Sprintf("%d:%02d", int(d.Minutes()), int(d.Seconds())%60)
	}
	return fmt.Sprintf("0:%02d", int(d.Seconds()))
}
