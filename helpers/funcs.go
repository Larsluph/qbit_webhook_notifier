package helpers

import (
	"fmt"
	"math"
	"time"
)

func FormatByteSize(bytes int64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB"}
	var i int
	var size = float64(bytes)
	for i = 0; size >= 1024 && i < len(sizes)-1; i++ {
		size /= 1024
	}
	return fmt.Sprintf("%.2f %s", size, sizes[i])
}

func FormatByteSpeed(bytes int64) string {
	return fmt.Sprintf("%s/s", FormatByteSize(bytes))
}

func RelativeTimeElapsed(timestamp1, timestamp2 int64) string {
	elapsed := time.Duration(math.Abs(float64(timestamp1-timestamp2))) * time.Second
	days := elapsed / (24 * time.Hour)
	elapsed -= days * 24 * time.Hour
	hours := elapsed / time.Hour
	elapsed -= hours * time.Hour
	minutes := elapsed / time.Minute
	elapsed -= minutes * time.Minute
	seconds := elapsed / time.Second

	if days > 0 {
		if days == 1 {
			return fmt.Sprintf("%d day, %d hours, %d minutes", days, hours, minutes)
		}
		return fmt.Sprintf("%d days, %d hours, %d minutes", days, hours, minutes)
	}

	if hours > 0 {
		if hours == 1 {
			return fmt.Sprintf("%d hour, %d minutes", hours, minutes)
		}
		return fmt.Sprintf("%d hours, %d minutes", hours, minutes)
	}

	if minutes > 0 {
		if minutes == 1 {
			return fmt.Sprintf("%d minute, %d seconds", minutes, seconds)
		}
		return fmt.Sprintf("%d minutes, %d seconds", minutes, seconds)
	}

	return fmt.Sprintf("%d seconds", seconds)
}
