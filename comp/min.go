package comp

type (
	Unsigned interface {
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
	}

	Signed interface {
		~int | ~int8 | ~int16 | ~int32 | ~int64
	}

	Integer interface {
		Signed | Unsigned
	}

	Float interface {
		~float32 | ~float64
	}

	Ordered interface {
		Integer | Float | ~string
	}
)

func Min[T Ordered](i, j T) T {
	if i > j {
		return j
	}

	return i
}

func Max[T Ordered](i, j T) T {
	if i < j {
		return j
	}

	return i
}
