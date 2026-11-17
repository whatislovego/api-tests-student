package utils

import (
	"testing"
	"time"
)

// WaitUntil периодически вызывает функцию check до возврата true или истечения timeout.
func WaitUntil(t *testing.T, timeout, interval time.Duration, message string, check func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if check() {
			return
		}
		time.Sleep(interval)
	}
	t.Fatalf("Таймаут ожидания (%s): %s", timeout, message)
}
