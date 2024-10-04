package helpers

import "time"

// FormatTimeToISO formats a time.Time to an ISO 8601 string.
func FormatTimeToISO(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ParseISOToTime parses an ISO 8601 string to a time.Time.
func ParseISOToTime(isoTime string) (time.Time, error) {
	return time.Parse(time.RFC3339, isoTime)
}

// GetCurrentTimestamp returns the current timestamp in Unix format.
func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}
