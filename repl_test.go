package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "hello world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "    sup   dude   ",
			expected: []string{"sup", "dude"},
		},
	}
	for _, c := range cases {
		output := cleanInput(c.input)
		if len(output) != len(c.expected) {
			t.Errorf("cleanInput is bOrKeN")
		}
		for i := range output {
			word := output[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("cleanInput is bOrKeN")
			}
		}
	}
}
