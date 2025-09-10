package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "   leading",
			expected: []string{"leading"},
		},
		{
			input:    "trailing   ",
			expected: []string{"trailing"},
		},
		{
			input:    "   surrounded   ",
			expected: []string{"surrounded"},
		},
		{
			input:    "multiple   spaces   here",
			expected: []string{"multiple", "spaces", "here"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("len(actual): %d does not match len(expected): %d", len(actual), len(c.expected))
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("word: %s is not matching expected: %s", word, expectedWord)
			}
		}
	}

}
