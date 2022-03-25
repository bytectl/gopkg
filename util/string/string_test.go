package string

import (
	"reflect"
	"testing"
)

func TestRemoveDuplicate(t *testing.T) {
	var tests = []struct {
		input []string
		want  []string
	}{
		{
			input: []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"},
			want:  []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"},
		},
		{
			input: []string{"aa", "bb", "aa", "", ""},
			want:  []string{"aa", "bb", ""},
		},
	}
	for _, test := range tests {
		got := RemoveDuplicate(test.input)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("RemoveDuplicate(%v) = %v; want %v", test.input, got, test.want)
		}
	}
}
