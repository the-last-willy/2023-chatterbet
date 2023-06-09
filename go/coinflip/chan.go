package coinflip

// ClearChannel empties a channel.
func ClearChannel[T any](c <-chan T) {
	for len(c) > 0 {
		_ = <-c
	}
}
