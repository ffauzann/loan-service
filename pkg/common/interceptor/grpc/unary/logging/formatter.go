package logging

import (
	"fmt"
	"time"
)

func formatMessage(start time.Time, method string) string {
	d := time.Since(start)
	msg := fmt.Sprintf("%s took %v", method, d)

	if ms := d.Milliseconds(); ms < 10000 {
		msg = fmt.Sprintf("%s took %vms", method, ms)
	}

	return msg
}
