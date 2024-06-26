package varfile

import (
	"testing"
)

func TestLexGroup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  []Var
	}{
		{
			name:  "one char",
			input: "[a]",
			want: []Var{
				{
					Type:  Group,
					Name:  "a",
					Value: "a",
				},
			},
		},
		{
			name:  "multiline",
			input: `
[a][a]
[b]
[c]
`,
			want: []Var{
				{
					Type:  Group,
					Name:  "a][a",
					Value: "a][a",
				},
				{
					Type:  Group,
					Name:  "b",
					Value: "b",
				},
				{
					Type:  Group,
					Name:  "c",
					Value: "c",
				},
			},
		},
		{
			name:  "multiline with incomplete line",
			input: `
[a][a]
[b
[c]
`,
			want: []Var{
				{
					Type:  Group,
					Name:  "a][a",
					Value: "a][a",
				},
				{
					Type:  Group,
					Name:  "c",
					Value: "c",
				},
			},
		},
		{
			name:  "one word",
			input: "[hello]",
			want: []Var{
				{
					Type:  Group,
					Name:  "hello",
					Value: "hello",
				},
			},
		},
		{
			name:  "nested group name",
			input: "[[hello]]",
			want: []Var{
				{
					Type:  Group,
					Name:  "[hello]",
					Value: "[hello]",
				},
			},
		},
		{
			name:  "space",
			input: "[hello world] [hello]",
			want: []Var{
				{
					Type:  Group,
					Name:  "hello",
					Value: "hello",
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			l := newLexer(test.name, test.input)
			vars, _ := l.ParseFile()
			if len(vars) != len(test.want) {
				t.Errorf("expected len(vars) %d, but got %d", len(test.want), len(vars))
			}

			for i := range vars {
				if vars[i] != test.want[i] {
					t.Errorf("expected %v, but %v", test.want[i], vars[i])
				}
			}
		})
	}
}


func TestLexComment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  []Var
	}{
		{
			name:  "one line",
			input: "# [comment]",
			want: []Var{},
		},
		{
			name:  "multiline",
			input: `
# [comment]
[a]
### [a]
`,
			want: []Var{
				{
					Type:  Group,
					Name:  "a",
					Value: "a",
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			l := newLexer(test.name, test.input)
			vars, _ := l.ParseFile()
			if len(vars) != len(test.want) {
				t.Errorf("expected len(vars) %d, but got %d", len(test.want), len(vars))
			}

			for i := range vars {
				if vars[i] != test.want[i] {
					t.Errorf("expected %v, but %v", test.want[i], vars[i])
				}
			}
		})
	}
}