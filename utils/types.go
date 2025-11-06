package utils

func StringPtr(s string) *string {
	return PtrOf(s)
}

func PtrOf[T any](val T) *T {
	return &val
}

func TypeToPtr[T any](ptr *T, defaultValue T) T {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

func ValOf[T any](ptr *T) T {
	if ptr == nil {
		var zeroVal T
		return zeroVal
	}
	return *ptr
}
