package repl

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "   ",
			expected: []string{},
		},
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "HELLO woRld",
			expected: []string{"hello", "world"},
		},
		{
			input:    "hello wORLD,    hello",
			expected: []string{"hello", "world,", "hello"},
		},
		{
			input:    "helloworld",
			expected: []string{"helloworld"},
		},
		{
			input:    " hello   ",
			expected: []string{"hello"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("lengths do not match: expected '%v', got '%v'", c.expected, actual)
			continue
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("actual word '%s' doesn't match expected '%s'", word, expectedWord)
			}
		}
	}
}
