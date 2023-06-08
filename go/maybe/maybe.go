package maybe

type Maybe[V any] struct {
	has   bool
	value V
}

func Just[V any](value V) Maybe[V] {
	return Maybe[V]{
		has:   true,
		value: value,
	}
}

func Nothing[V any]() Maybe[V] {
	return Maybe[V]{
		has: false,
	}
}

func (o *Maybe[V]) Value() (value V, has bool) {
	return o.value, o.has
}
