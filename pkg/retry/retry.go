package retry

import "time"

// Do starts trying to exec the given fn. If it returns an error, it will wait the given amount of time and retry, for `maxRetries` times
func Do(fn func() error, wait time.Duration, maxRetries int) error {
	for i := 0; i < maxRetries; i++ {
		err := fn()
		if err == nil {
			return nil
		}
		time.Sleep(wait)
	}
	return nil
}
