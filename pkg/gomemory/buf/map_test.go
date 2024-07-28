package buf_test

import (
	"testing"

	"github.com/jejikeh/invaders/pkg/gomemory/buf"
)

type A struct {
	a, b int
	c    string
}

func TestBufMap(t *testing.T) {
	tests := []struct {
		name   string
		count  int
		load   func(*testing.T, *buf.Map[int, A])
		assert func(*testing.T, *buf.Map[int, A])
	}{
		{
			name:  "single int",
			count: 1,
			load: func(t *testing.T, b *buf.Map[int, A]) {
				t.Helper()

				a := b.Get(1)
				*a = A{a: 1, b: 2, c: "foo"}
			},
			assert: func(t *testing.T, b *buf.Map[int, A]) {
				t.Helper()

				a := b.Get(1)
				if a.a != 1 || a.b != 2 || a.c != "foo" {
					t.Fail()
				}
			},
		},
		{
			name:  "two int array",
			count: 2,
			load: func(t *testing.T, b *buf.Map[int, A]) {
				t.Helper()

				a := b.Get(1)
				*a = A{a: 1, b: 2, c: "foo"}
				a = b.Get(2)
				*a = A{a: 3, b: 4, c: "bar"}
			},
			assert: func(t *testing.T, b *buf.Map[int, A]) {
				t.Helper()

				a := b.Get(1)
				if a.a != 1 || a.b != 2 || a.c != "foo" {
					t.Fail()
				}

				a = b.Get(2)
				if a.a != 3 || a.b != 4 || a.c != "bar" {
					t.Fail()
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			b := buf.NewMap[int, A](test.count)
			test.load(t, b)
			test.assert(t, b)
		})
	}
}
