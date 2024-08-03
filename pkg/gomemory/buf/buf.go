package buf

import "errors"

const BufChunkSize = 1024

var ErrBufOverflow = errors.New("buffer overflow")

type Buf[T any] struct {
	buf [][BufChunkSize]T

	index int
	bufId int
}

func New[T any](count int) *Buf[T] {
	return &Buf[T]{
		buf: make([][BufChunkSize]T, count/BufChunkSize+1),
	}
}

func (b *Buf[T]) New() *T {
	if b.index >= len(b.buf[b.bufId]) {
		b.bufId++
		b.index = 0
	}

	if b.bufId >= len(b.buf) {
		panic(ErrBufOverflow)
	}

	defer func() {
		b.index++
	}()

	return &b.buf[b.bufId][b.index]
}

func (b *Buf[T]) Len() int {
	return b.index
}

func (b *Buf[T]) Reset() {
	b.bufId = 0
	b.index = 0
}

func (b *Buf[T]) Clear() {
	b.bufId = 0
	b.index = 0

	for i := range b.buf {
		b.buf[i] = [BufChunkSize]T{}
	}
}

func (b *Buf[T]) Get(idx int) *T {
	return &b.buf[b.bufId][idx]
}
