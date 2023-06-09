package maybe

// Maybe instances might or might not contain a value of a given type.
type Maybe[T any] struct {
	has   bool
	value T
}

// Just returns a `Maybe` containing `value`
func Just[T any](value T) Maybe[T] {
	return Maybe[T]{
		has:   true,
		value: value,
	}
}

// Nothing returns a `Maybe` containing no value.
func Nothing[T any]() Maybe[T] {
	return Maybe[T]{
		has: false,
	}
}

// Value returns the value of maybe and whether that value is valid.
func (o *Maybe[T]) Value() (value T, has bool) {
	return o.value, o.has
}
