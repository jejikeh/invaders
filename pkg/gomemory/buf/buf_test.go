package buf_test

import (
	"testing"

	"github.com/jejikeh/invaders/pkg/gomemory/buf"
)

type testData struct {
	a int
	b int
}

func TestBuf(t *testing.T) {
	tests := []struct {
		name  string
		count int
		data  testData
	}{
		{
			name:  "single int",
			count: 1,
			data: testData{
				a: 1,
				b: 2,
			},
		},
		{
			name:  "two int array",
			count: 2,
			data: testData{
				a: 1,
				b: 2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			b := buf.New[testData](tt.count)

			for range tt.count {
				a := b.New()
				*a = tt.data
			}

			for i := range tt.count {
				got := b.Get(i)
				if got.a != tt.data.a || got.b != tt.data.b {
					t.Errorf("buf.New() = %v, want %v", got, tt.data)
				}
			}
		})
	}
}

func TestBufReset(t *testing.T) {
	tests := []struct {
		name  string
		count int
		data  testData
	}{
		{
			name:  "single int",
			count: 1,
			data: testData{
				a: 1,
				b: 2,
			},
		},
		{
			name:  "two int array",
			count: 2,
			data: testData{
				a: 1,
				b: 2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			b := buf.New[testData](tt.count)

			for range tt.count {
				a := b.New()
				*a = tt.data
			}

			b.Reset()

			for range tt.count {
				a := b.New()
				*a = tt.data
			}

			length := b.Len()
			if length != tt.count {
				t.Errorf("buf.New() = %v, want %v", length, tt.count)
			}

			for i := range tt.count {
				got := b.Get(i)
				if got.a != tt.data.a || got.b != tt.data.b {
					t.Errorf("buf.New() = %v, want %v", got, tt.data)
				}
			}
		})
	}
}

func TestBufClear(t *testing.T) {
	tests := []struct {
		name  string
		count int
		data  testData
	}{
		{
			name:  "single int",
			count: 1,
			data: testData{
				a: 1,
				b: 2,
			},
		},
		{
			name:  "two int array",
			count: 2,
			data: testData{
				a: 1,
				b: 2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			b := buf.New[testData](tt.count)

			for range tt.count {
				a := b.New()
				*a = tt.data
			}

			b.Clear()

			length := b.Len()
			if length != 0 {
				t.Errorf("buf.New() = %v, want %v", length, 0)
			}

			for i := range tt.count {
				got := b.Get(i)
				if got.a != 0 || got.b != 0 {
					t.Errorf("buf.New() = %v, want %v", got, testData{})
				}
			}
		})
	}
}
