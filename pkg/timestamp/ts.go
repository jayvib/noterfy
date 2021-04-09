package timestamp

import (
	"noteapp/pkg/ptrconv"
	"time"
)

// GenerateTimestamp return the current truncated(in second) time
// pointer in UTC.
func GenerateTimestamp() *time.Time {
	return ptrconv.TimePointer(time.Now().UTC().Truncate(time.Second))
}
