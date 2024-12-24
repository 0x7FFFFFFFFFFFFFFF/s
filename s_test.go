package s

import (
	"reflect"
	"testing"
)

func TestUnindent(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`
		hello world
		`, "hello world"},
		{`
		
		hello world
		`, "\nhello world"},
		{`    
		
	hello world
		`, "\t\nhello world"},
		{"   spaces   ", "spaces   "},
		{"", ""},
		{" ", ""},
		{" \t\t", ""},
	}

	for _, test := range tests {
		result := Unindent(test.input)
		if result != test.expected {
			t.Errorf("Unindent(%q) = %q; want %q",
				test.input, result, test.expected)
		}
	}
}

func TestTrim(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"   spaces   ", "spaces"},
		{"\t\ttabs\t\t", "tabs"},
		{"\r\t\ttabs\t\t\n\r", "tabs"},
		{"\r\n\t\t\r\n\r\n", ""},
		{WhiteSpace + "hello world" + WhiteSpace, "hello world"},
	}

	for _, test := range tests {
		result := Trim(test.input)
		if result != test.expected {
			t.Errorf("Trim(%q) = %q; want %q",
				test.input, result, test.expected)
		}
	}
}

func TestTrimLeft(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"   spaces   ", "spaces   "},
		{"\t\ttabs\t\t", "tabs\t\t"},
		{"\r\t\ttabs\t\t\n\r", "tabs\t\t\n\r"},
		{"\r\n\t\t\r\n\r\n", ""},
		{WhiteSpace + "hello world" + WhiteSpace, "hello world" + WhiteSpace},
	}

	for _, test := range tests {
		result := TrimLeft(test.input)
		if result != test.expected {
			t.Errorf("TrimLeft(%q) = %q; want %q",
				test.input, result, test.expected)
		}
	}
}

func TestTrimRight(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"   spaces   ", "   spaces"},
		{"\t\ttabs\t\t", "\t\ttabs"},
		{"\r\t\ttabs\t\t\n\r", "\r\t\ttabs"},
		{"\r\n\t\t\r\n\r\n", ""},
		{WhiteSpace + "hello world" + WhiteSpace, WhiteSpace + "hello world"},
	}

	for _, test := range tests {
		result := TrimRight(test.input)
		if result != test.expected {
			t.Errorf("TrimRight(%q) = %q; want %q",
				test.input, result, test.expected)
		}
	}
}

func TestGetIndentString(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"\r\n\t\t\r\n\r\n", 0},
		{"   spaces   ", 3},
		{"   spaces\nabc", 0},
		{"\t\ttabs\t\t", 2},
		{"\r\t\ttabs\t\t\n\r", 0},
		{"  \r\n\t\t  \n  abcd\r\n  ", 2},
		{"  \r\n\t\t  \n  abcd\r\n", 0},
		{"  \r\n\t\t  \n  abcd\r", 2},
		{"  \r\n\t\t  \n  abcd\r\n  cdef", 2},
	}

	for _, test := range tests {
		result := GetIndentString(test.input)
		if result != test.expected {
			t.Errorf("GetIndentString(%q) = %d; want %d",
				test.input, result, test.expected)
		}
	}
}

func TestGetIndentStringArray(t *testing.T) {
	tests := []struct {
		input    []string
		expected int
	}{
		{[]string{"\r\n", "\t\t\r\n", "\r\n"}, 0},
		{[]string{"   spaces   "}, 3},
		{[]string{"   spaces\n", "abc"}, 0},
		{[]string{"\t\ttabs\t\t"}, 2},
		{[]string{"\r\t\ttabs\t\t\n\r"}, 0},
		{[]string{"  \r\n", "\t\t  \n", "  abcd\r\n", "  "}, 2},
		{[]string{"  \r\n", "\t\t  \n", "  abcd\r\n"}, 2},
		{[]string{"  \r\n", "\t\t  \n", "  abcd\r"}, 2},
		{[]string{"  \r\n", "\t\t  \n", "  abcd\r\n", "  cdef"}, 2},
	}

	for _, test := range tests {
		result := GetIndentStringArray(test.input)
		if result != test.expected {
			t.Errorf("GetIndentStringArray(%q) = %d; want %d",
				test.input, result, test.expected)
		}
	}
}

func TestToWindowsLineEnding(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"\r\n", "\r\n"},
		{"\n", "\r\n"},
		{"\r", "\r"},
		{"\r\n\n\r", "\r\n\r\n\r"},
		{"\n\r\n", "\r\n\r\n"},
		{"\r\n\r\n", "\r\n\r\n"},
		{"\n\n", "\r\n\r\n"},
		{"\r\r", "\r\r"},
		{"", ""},
	}

	for _, test := range tests {
		result := ToWindowsLineEnding(test.input)
		if result != test.expected {
			t.Errorf("ToWindowsLineEnding(%q) = %q; want %q",
				test.input, result, test.expected)
		}
	}
}

func TestSucc(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0", "1"},
		{"00", "01"},
		{"1", "2"},
		{"a", "b"},
		{"z", "aa"},

		// Basic alphanumeric cases
		{"THX1138", "THX1139"},
		{"<<koala>>", "<<koalb>>"},
		{"***", "**+"},

		// Digit cases with carrying
		{"00", "01"},
		{"09", "10"},
		{"99", "100"},

		// Lowercase letter cases with carrying
		{"aa", "ab"},
		{"az", "ba"},
		{"zz", "aaa"},

		// Uppercase letter cases with carrying
		{"AA", "AB"},
		{"AZ", "BA"},
		{"ZZ", "AAA"},

		// Mixed cases with carrying
		{"zz99zz99", "aaa00aa00"},
		{"99zz99zz", "100aa00aa"},

		// Empty string
		{"", ""},

		// Single character cases
		{"a", "b"},
		{"z", "aa"},
		{"Z", "AA"},
		{"9", "10"},

		// Additional edge cases
		{"1999zzz", "2000aaa"},
		{"ZZZ9999", "AAAA0000"},
		{"abcdefghijklmnopqrstuvwxyz", "abcdefghijklmnopqrstuvwxza"},
	}

	for _, test := range tests {
		result := Succ(test.input)
		if result != test.expected {
			t.Errorf("Succ(%q) = %q; want %q",
				test.input, result, test.expected)
		}
	}
}

func TestHighlight(t *testing.T) {
	tests := []struct {
		input    string
		pattern  string
		left     string
		right    string
		expected string
	}{
		{
			input:    "test string test",
			pattern:  "test",
			left:     "<",
			right:    ">",
			expected: "<test> string <test>",
		},
		{
			input:    "hello world",
			pattern:  "o",
			left:     "(",
			right:    ")",
			expected: "hell(o) w(o)rld",
		},
		{
			input:    "no matches here",
			pattern:  "xyz",
			left:     "<",
			right:    ">",
			expected: "no matches here",
		},
		{
			input:    "overlapping pattern",
			pattern:  "pattern",
			left:     "«",
			right:    "»",
			expected: "overlapping «pattern»",
		},
		{
			input:    "",
			pattern:  "test",
			left:     "<",
			right:    ">",
			expected: "",
		},
		{
			input:    "multiple   spaces",
			pattern:  "\\s+",
			left:     "[",
			right:    "]",
			expected: "multiple[   ]spaces",
		},
		{
			input:    "  line 1  \r\n  line 2  \r\n  line \t33  ",
			pattern:  `line\s+(\d+)`,
			left:     "[",
			right:    "]",
			expected: "  [line 1]  \r\n  [line 2]  \r\n  [line \t33]  ",
		},
	}

	for _, test := range tests {
		result := Highlight(test.input, test.pattern, test.left, test.right)
		if result != test.expected {
			t.Errorf("HighlightString(%q, %q, %q, %q) = %q; want %q",
				test.input, test.pattern, test.left, test.right,
				result, test.expected)
		}
	}
}

func TestIndent(t *testing.T) {
	tests := []struct {
		input    string
		indent   string
		expected string
	}{
		{
			input:    "line1\nline2\nline3",
			indent:   "  ",
			expected: "  line1\n  line2\n  line3",
		},
		{
			input:    "line1\r\nline2\r\nline3",
			indent:   "\t",
			expected: "\tline1\r\n\tline2\r\n\tline3",
		},
		{
			input:    "line1\rline2\rline3",
			indent:   ">>",
			expected: ">>line1\rline2\rline3",
		},
		{
			input:    "\nline1\n\nline2\n",
			indent:   "  ",
			expected: "  \n  line1\n  \n  line2\n  ",
		},
		{
			input:    "",
			indent:   "  ",
			expected: "  ",
		},
		{
			input:    "single",
			indent:   "-->",
			expected: "-->single",
		},
		{
			input:    "line1\n  line2\n    line3",
			indent:   "  ",
			expected: "  line1\n    line2\n      line3",
		},
	}

	for _, test := range tests {
		result := Indent(test.input, test.indent)
		if result != test.expected {
			t.Errorf("\nInput:\n%q\nIndent:\n%q\nGot:\n%q\nWant:\n%q",
				test.input, test.indent, result, test.expected)
		}
	}
}

func TestLen(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"böt", 3},
		{"hello", 5},
		{"世界", 2},
		{"café", 4},
		{"", 0},
		{"π", 1},
		{"🍕", 1},
		{"한글", 2},
		{"böt世界", 5},
		{"a\u0308", 2}, // 'a' with umlaut combining character
		{"é", 1},       // single character é
		{"e\u0301", 2}, // 'e' with acute accent combining character
	}

	for _, test := range tests {
		result := Len(test.input)
		if result != test.expected {
			t.Errorf("Len(%q) = %d; want %d",
				test.input, result, test.expected)
		}
	}
}

func TestEachChar(t *testing.T) {
	tests := []struct {
		input    string
		expected []struct {
			char  string
			index int
		}
	}{
		{
			"hello",
			[]struct {
				char  string
				index int
			}{
				{"h", 0},
				{"e", 1},
				{"l", 2},
				{"l", 3},
				{"o", 4},
			},
		},
		{
			"世界",
			[]struct {
				char  string
				index int
			}{
				{"世", 0},
				{"界", 1},
			},
		},
		{
			"café",
			[]struct {
				char  string
				index int
			}{
				{"c", 0},
				{"a", 1},
				{"f", 2},
				{"é", 3},
			},
		},
		{
			"🍕π",
			[]struct {
				char  string
				index int
			}{
				{"🍕", 0},
				{"π", 1},
			},
		},
		{
			"",
			nil,
		},
	}

	for _, test := range tests {
		var result []struct {
			char  string
			index int
		}

		EachChar(test.input, func(char string, index int) {
			result = append(result, struct {
				char  string
				index int
			}{char, index})
		})

		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("EachChar(%q) = %v; want %v",
				test.input, result, test.expected)
		}
	}

}

func TestMultiple(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		affixes     []string
		expectedPre string
		expectedApp string
	}{
		{
			name:        "multiple basic",
			input:       "world",
			affixes:     []string{"hello ", "hey ", "hi "},
			expectedPre: "hi hey hello world",
			expectedApp: "worldhello hey hi ",
		},
		{
			name:        "single affix",
			input:       "test",
			affixes:     []string{"one"},
			expectedPre: "onetest",
			expectedApp: "testone",
		},
		{
			name:        "empty input",
			input:       "",
			affixes:     []string{"1", "2", "3"},
			expectedPre: "321",
			expectedApp: "123",
		},
		{
			name:        "empty affixes",
			input:       "test",
			affixes:     []string{},
			expectedPre: "test",
			expectedApp: "test",
		},
		{
			name:        "unicode multiple",
			input:       "界",
			affixes:     []string{"世", "家"},
			expectedPre: "家世界",
			expectedApp: "界世家",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resultPre := Prepend(test.input, test.affixes...)
			if resultPre != test.expectedPre {
				t.Errorf("Prepend(%q, %v) = %q; want %q",
					test.input, test.affixes, resultPre, test.expectedPre)
			}

			resultApp := Append(test.input, test.affixes...)
			if resultApp != test.expectedApp {
				t.Errorf("Append(%q, %v) = %q; want %q",
					test.input, test.affixes, resultApp, test.expectedApp)
			}
		})
	}
}
func TestIsMatch(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		pattern string
		want    bool
	}{
		// Full match tests (with ^ and $)
		{
			name:    "exact match with anchors",
			s:       "hello",
			pattern: "^hello$",
			want:    true,
		},
		{
			name:    "non-match with anchors",
			s:       "hello world",
			pattern: "^hello$",
			want:    false,
		},

		// Partial match tests (without ^ and $)
		{
			name:    "partial match at start",
			s:       "hello world",
			pattern: "hello",
			want:    true,
		},
		{
			name:    "partial match in middle",
			s:       "say hello world",
			pattern: "hello",
			want:    true,
		},
		{
			name:    "partial match at end",
			s:       "say hello",
			pattern: "hello",
			want:    true,
		},

		// Mixed anchor tests
		{
			name:    "start anchor only",
			s:       "hello world",
			pattern: "^hello",
			want:    true,
		},
		{
			name:    "end anchor only",
			s:       "say hello",
			pattern: "hello$",
			want:    true,
		},

		// Special regex features
		{
			name:    "wildcard match",
			s:       "hello",
			pattern: "h.*o",
			want:    true,
		},
		{
			name:    "character class",
			s:       "hello",
			pattern: "[a-z]+",
			want:    true,
		},
		{
			name:    "alternation",
			s:       "hello",
			pattern: "hello|world",
			want:    true,
		},

		// Unicode tests
		{
			name:    "unicode match",
			s:       "你好世界",
			pattern: "世界",
			want:    true,
		},
		{
			name:    "unicode with anchors",
			s:       "你好世界",
			pattern: "^你好",
			want:    true,
		},

		// Empty string tests
		{
			name:    "empty string match",
			s:       "",
			pattern: "^$",
			want:    true,
		},
		{
			name:    "empty string non-match",
			s:       "",
			pattern: ".",
			want:    false,
		},

		// Invalid regex pattern
		{
			name:    "invalid regex",
			s:       "hello",
			pattern: "[",
			want:    false,
		},

		// Quantifier tests
		{
			name:    "zero or more",
			s:       "aaa",
			pattern: "a*",
			want:    true,
		},
		{
			name:    "one or more",
			s:       "aaa",
			pattern: "a+",
			want:    true,
		},

		// Word boundary tests
		{
			name:    "word boundary",
			s:       "hello world",
			pattern: "\\bhello\\b",
			want:    true,
		},

		// Case sensitivity tests
		{
			name:    "case sensitive",
			s:       "Hello",
			pattern: "hello",
			want:    false,
		},
		{
			name:    "case insensitive",
			s:       "Hello",
			pattern: "(?i)hello",
			want:    true,
		},

		// Special characters in input
		{
			name:    "special chars in input",
			s:       "hello.world",
			pattern: "hello\\.world",
			want:    true,
		},

		// Capturing groups
		{
			name:    "capturing group",
			s:       "hello123world",
			pattern: "hello(\\d+)world",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsMatch(tt.s, tt.pattern)
			if got != tt.want {
				t.Errorf("IsMatch(%q, %q) = %v; want %v",
					tt.s, tt.pattern, got, tt.want)
			}
		})
	}
}

func TestGrep(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		pattern string
		want    []string
	}{
		{
			name:    "simple word match",
			s:       "hello world hello universe",
			pattern: "hello",
			want:    []string{"hello", "hello"},
		},
		{
			name:    "with capture group",
			s:       "hello world hello universe",
			pattern: "(hello)",
			want:    []string{"hello", "hello"},
		},
		{
			name:    "numbers",
			s:       "abc123def456ghi789",
			pattern: "\\d+",
			want:    []string{"123", "456", "789"},
		},
		{
			name:    "words",
			s:       "The quick brown fox",
			pattern: "\\w+",
			want:    []string{"The", "quick", "brown", "fox"},
		},
		{
			name:    "email addresses",
			s:       "Contact us at: test@example.com or support@example.com",
			pattern: "[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}",
			want:    []string{"test@example.com", "support@example.com"},
		},
		{
			name:    "specific capture group",
			s:       "name: John, age: 30, name: Jane, age: 25",
			pattern: "name: ([^,]+)",
			want:    []string{"name: John", "name: Jane"},
		},
		{
			name:    "unicode",
			s:       "你好世界你好宇宙",
			pattern: "你好",
			want:    []string{"你好", "你好"},
		},
		{
			name:    "empty input",
			s:       "",
			pattern: "\\w+",
			want:    []string{},
		},
		{
			name:    "no matches",
			s:       "hello world",
			pattern: "xyz",
			want:    []string{},
		},
		{
			name:    "invalid regex",
			s:       "hello world",
			pattern: "[",
			want:    []string{},
		},
		{
			name:    "overlapping matches",
			s:       "aaaaa",
			pattern: "aa",
			want:    []string{"aa", "aa"},
		},
		{
			name:    "multiple capture groups",
			s:       "name: John, age: 30, name: Jane, age: 25",
			pattern: "(name): ([^,]+)",
			want:    []string{"name: John", "name: Jane"},
		},
		{
			name:    "case insensitive",
			s:       "Hello HELLO hello",
			pattern: "(?i)hello",
			want:    []string{"Hello", "HELLO", "hello"},
		},
		{
			name:    "word boundaries",
			s:       "hello hello123 hello",
			pattern: "\\bhello\\b",
			want:    []string{"hello", "hello"},
		},
		{
			name:    "multiline",
			s:       "hello\nworld\nhello",
			pattern: "(?m)^hello$",
			want:    []string{"hello", "hello"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Grep(tt.s, tt.pattern)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Grep(%q, %q) = %v; want %v",
					tt.s, tt.pattern, got, tt.want)
			}
		})
	}
}