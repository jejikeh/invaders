package gomath

func IsPowerOfTwo[T int | uint | uintptr](x T) bool {
	return x != 0 && x&(x-1) == 0
}
