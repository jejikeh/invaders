package rtype

import (
	"testing"
)

func TestGetITab(t *testing.T) {
	tests := []struct {
		name    string
		prepare func() []ITab
		check   func(*testing.T, []ITab)
	}{
		{
			name: "simple ints",
			prepare: func() []ITab {
				return []ITab{
					GetITab(1),
					GetITab(2),
					GetITab(3),
					GetITab(4),
				}
			},
			check: func(t *testing.T, res []ITab) {
				t.Helper()

				if !(res[0] == res[1] && res[1] == res[2] && res[2] == res[3]) {
					t.Errorf("got %v, want %v", res, res)
				}
			},
		},
		{
			name: "simple pointers",
			prepare: func() []ITab {
				return []ITab{
					GetITab(&struct{}{}),
					GetITab(&struct{}{}),
					GetITab(&struct{}{}),
					GetITab(&struct{}{}),
				}
			},
			check: func(t *testing.T, res []ITab) {
				t.Helper()

				if !(res[0] == res[1] && res[1] == res[2] && res[2] == res[3]) {
					t.Errorf("got %v, want %v", res, res)
				}
			},
		},
		{
			name: "simple arrays",
			prepare: func() []ITab {
				return []ITab{
					GetITab([]int{1, 2, 3, 4}),
					GetITab([]int{1, 2, 3}),
					GetITab([]int{1, 2}),
					GetITab([]int{1}),
				}
			},
			check: func(t *testing.T, res []ITab) {
				t.Helper()

				if !(res[0] == res[1] && res[1] == res[2] && res[2] == res[3]) {
					t.Errorf("got %v, want %v", res, res)
				}
			},
		},
		{
			name: "simple slices",
			prepare: func() []ITab {
				return []ITab{
					GetITab([]int{1, 2, 3, 4}),
					GetITab([]int{1, 2, 3}),
					GetITab([]int{1, 2}),
					GetITab([]int{1}),
				}
			},
			check: func(t *testing.T, res []ITab) {
				t.Helper()

				if !(res[0] == res[1] && res[1] == res[2] && res[2] == res[3]) {
					t.Errorf("got %v, want %v", res, res)
				}
			},
		},
		{
			name: "different types",
			prepare: func() []ITab {
				return []ITab{
					GetITab([]int32{1, 2, 3, 4}),
					GetITab([]int16{1, 2, 3}),
					GetITab([]uint32{1, 2}),
					GetITab([]uint16{1}),
				}
			},
			check: func(t *testing.T, res []ITab) {
				t.Helper()

				if !(res[0] != res[1] && res[1] != res[2] && res[2] != res[3]) {
					t.Errorf("got %v, want %v", res, res)
				}
			},
		},
	}

	for _, tests := range tests {
		t.Run(tests.name, func(t *testing.T) {
			t.Parallel()

			res := tests.prepare()
			tests.check(t, res)
		})
	}
}
