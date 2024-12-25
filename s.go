package s

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

const WhiteSpace = "\t\n\r "

func Trim(s string) string {
	return strings.Trim(s, WhiteSpace)
}

func TrimLeft(s string) string {
	return strings.TrimLeft(s, WhiteSpace)
}

func TrimRight(s string) string {
	return strings.TrimRight(s, WhiteSpace)
}

func Strip(s string) string {
	return Trim(s)
}

func LStrip(s string) string {
	return TrimLeft(s)
}

func RStrip(s string) string {
	return TrimRight(s)
}

func GetIndentString(s string) int {
	// Normalize all line endings to \n
	s = ToLinuxLineEnding(s)

	lines := strings.Split(s, "\n")

	minIndentSize := math.MaxInt
	for _, line := range lines {
		line = strings.TrimRight(line, "\r\n")
		indentSize := len(line) - len(strings.TrimLeft(line, " \t"))

		if indentSize < minIndentSize {
			minIndentSize = indentSize
		}
	}
	return minIndentSize
}

func GetIndentStringArray(lines []string) int {
	minIndentSize := math.MaxInt
	for _, line := range lines {
		line = strings.TrimRight(line, "\r\n")
		indentSize := len(line) - len(strings.TrimLeft(line, " \t"))

		if indentSize < minIndentSize {
			minIndentSize = indentSize
		}
	}
	return minIndentSize
}

func GetLineEnding(s string) string {
	if strings.Contains(s, "\r\n") {
		return "\r\n"
	}
	return "\n"
}

func ToLinuxLineEnding(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

func ToWindowsLineEnding(s string) string {
	reg := regexp.MustCompile(`\r?\n`)
	return reg.ReplaceAllString(s, "\r\n")
}

// - A string with the common leading whitespace removed from each line and original line ending style restored.
func Unindent(s string) string {
	// Detect original line ending style
	lineEnding := GetLineEnding(s)

	// Normalize all line endings to \n
	s = ToLinuxLineEnding(s)

	// removeIndentation removes n characters from the front of each line in lines.
	removeIndentation := func(lines []string, n int) []string {
		for i, line := range lines {
			line = strings.TrimRight(line, "\r\n")
			if len(line) >= n {
				lines[i] = line[n:]
			}
		}
		return lines
	}

	// Split on \n since we've already normalized line endings
	lines := strings.Split(s, "\n")
	if len(lines) > 0 && Trim(lines[0]) == "" {
		lines = lines[1:]
	}

	if len(lines) > 0 && Trim(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}

	indent := GetIndentStringArray(lines)
	lines = removeIndentation(lines, indent)

	return strings.Join(lines, lineEnding)
}

func Succ(s string) string {
	if len(s) == 0 {
		return s
	}

	runes := []rune(s)
	pos := len(runes) - 1

	// Find the rightmost character to increment
	// For alphanumeric strings, find rightmost alphanumeric
	// For non-alphanumeric strings, use rightmost character
	lastAlphaNumPos := -1
	for i := len(runes) - 1; i >= 0; i-- {
		if unicode.IsLetter(runes[i]) || unicode.IsDigit(runes[i]) {
			lastAlphaNumPos = i
			break
		}
	}

	if lastAlphaNumPos == -1 {
		// No alphanumeric found, increment the last character
		runes[len(runes)-1]++
		return string(runes)
	}

	// Handle the increment starting from the rightmost alphanumeric
	pos = lastAlphaNumPos
	carry := true

	for pos >= 0 && carry {
		r := runes[pos]

		if unicode.IsDigit(r) {
			if r == '9' {
				runes[pos] = '0'
				carry = true
			} else {
				runes[pos]++
				carry = false
			}
		} else if unicode.IsLetter(r) {
			if unicode.IsUpper(r) {
				if r == 'Z' {
					runes[pos] = 'A'
					carry = true
				} else {
					runes[pos]++
					carry = false
				}
			} else {
				if r == 'z' {
					runes[pos] = 'a'
					carry = true
				} else {
					runes[pos]++
					carry = false
				}
			}
		} else {
			// For non-alphanumeric characters
			runes[pos]++
			carry = false
		}
		pos--
	}

	// Handle carrying beyond the leftmost position
	if carry {
		prefix := ""
		if pos < 0 {
			firstChar := runes[0]
			if unicode.IsDigit(firstChar) {
				prefix = "1"
			} else if unicode.IsLetter(firstChar) {
				if unicode.IsUpper(firstChar) {
					prefix = "A"
				} else {
					prefix = "a"
				}
			}
		}
		return prefix + string(runes)
	}

	return string(runes)
}

// Highlight highlights all occurrences of a pattern in a given string by surrounding them with specified left and right markers.
//
// Parameters:
//   - s: The input string in which to search for the pattern.
//   - pattern: The regular expression pattern to search for in the input string.
//   - left: The string to insert to the left of each match.
//   - right: The string to insert to the right of each match.
//
// Returns:
//
//	A new string with all occurrences of the pattern surrounded by the left and right markers. If the pattern is not found or if there is an error compiling the pattern, the original string is returned.
func Highlight(s string, pattern string, left string, right string) string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return s
	}

	matches := re.FindAllIndex([]byte(s), -1)
	if len(matches) == 0 {
		return s
	}

	var builder strings.Builder
	lastPos := 0

	// Pre-calculate the total size to avoid reallocations
	totalSize := len(s) + (len(left)+len(right))*len(matches)
	builder.Grow(totalSize)

	// Process each match
	for _, match := range matches {
		start, end := match[0], match[1]

		// Write the text between the last match and this match
		builder.WriteString(s[lastPos:start])

		// Write the highlight markers and the matched text
		builder.WriteString(left)
		builder.WriteString(s[start:end])
		builder.WriteString(right)

		lastPos = end
	}

	// Write any remaining text after the last match
	if lastPos < len(s) {
		builder.WriteString(s[lastPos:])
	}

	return builder.String()
}

func Indent(s string, indent string) string {
	// Detect original line ending style
	lineEnding := GetLineEnding(s)

	// Normalize all line endings to \n
	s = ToLinuxLineEnding(s)

	// Split into lines
	lines := strings.Split(s, "\n")

	// Process each line
	for i, line := range lines {
		lines[i] = indent + line
	}

	// Join lines with original line ending
	return strings.Join(lines, lineEnding)
}

func LenRune(s string) int {
	return len([]rune(s))
}

func LengthRune(s string) int {
	return LenRune(s)
}

func LengthByte(s string) int {
	return len(s)
}

func LenByte(s string) int {
	return len(s)
}

func EachChar(s string, callback func(char string, index int)) {
	runes := []rune(s)
	for i, r := range runes {
		callback(string(r), i)
	}
}

// Append adds multiple strings to the end of s
func Append(s string, suffixes ...string) string {
	var builder strings.Builder
	builder.Grow(len(s) + sumLength(suffixes))

	builder.WriteString(s)
	for _, suffix := range suffixes {
		builder.WriteString(suffix)
	}

	return builder.String()
}

// Prepend adds multiple strings to the start of s
func Prepend(s string, prefixes ...string) string {
	var builder strings.Builder
	builder.Grow(len(s) + sumLength(prefixes))

	for i := len(prefixes) - 1; i >= 0; i-- {
		builder.WriteString(prefixes[i])
	}
	builder.WriteString(s)

	return builder.String()
}

// Helper function to calculate total length of strings
func sumLength(strings []string) int {
	total := 0
	for _, s := range strings {
		total += len(s)
	}
	return total
}

func IsMatch(s string, regex string) bool {
	// If regex doesn't start with ^, allow partial matches from start
	if !strings.HasPrefix(regex, "^") {
		regex = ".*" + regex
	}

	// If regex doesn't end with $, allow partial matches at end
	if !strings.HasSuffix(regex, "$") {
		regex = regex + ".*"
	}

	re, err := regexp.Compile(regex)
	if err != nil {
		return false
	}

	return re.MatchString(s)
}

func Grep(s string, pattern string) []string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return []string{}
	}

	// Find all matches
	matches := re.FindAllStringSubmatch(s, -1)
	if matches == nil {
		return []string{}
	}

	// Extract the full matches
	result := make([]string, 0, len(matches))
	for _, match := range matches {
		result = append(result, match[0])
	}

	return result
}

func GetGroup(s, pattern, group string) string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return ""
	}

	match := re.FindStringSubmatch(s)
	if match == nil {
		return ""
	}

	// If group is a number
	if groupNum, err := strconv.Atoi(group); err == nil {
		if groupNum < 0 || groupNum >= len(match) {
			return ""
		}
		return match[groupNum]
	}

	// If group is a name
	names := re.SubexpNames()
	for i, name := range names {
		if name == group && i < len(match) {
			return match[i]
		}
	}

	return ""
}

// Repeat returns s repeated count times (concatenated).
// e.g. Repeat("XY", 3) => "XYXYXY".
func Repeat(s string, count int) string {
	if count <= 0 {
		return ""
	}

	var builder strings.Builder
	builder.Grow(len(s) * count)

	for i := 0; i < count; i++ {
		builder.WriteString(s)
	}

	return builder.String()
}

// LeftPad: repeat padStr on the left. If final exceeds length, trim from the right.
func LeftPad(s string, padStr string, length int) string {
	runes := []rune(s)
	if len(runes) >= length || padStr == "" {
		return s
	}
	needed := length - len(runes)
	prepend := Repeat(padStr, needed)
	// Grab "needed" runes from the beginning of prepend
	result := []rune(prepend)[:needed]
	result = append(result, runes...)
	return string(result)
}

// RightPad: repeat padStr on the right.
func RightPad(s string, padStr string, length int) string {
	runes := []rune(s)
	if len(runes) >= length || padStr == "" {
		return s
	}
	needed := length - len(runes)
	appendStr := Repeat(padStr, needed)
	// Grab "needed" runes from the end of appendStr
	result := append(runes, []rune(appendStr)[:needed]...)
	return string(result)
}

// Pad: two-sided padding in a stepwise manner (left first, then right, then left, etc.).
func Pad(s string, padStr string, length int) string {
	runes := []rune(s)
	s_len := len(runes)
	if len(runes) >= length || padStr == "" {
		return s
	}
	pad_len := len([]rune(padStr))
	left_len := 0
	right_len := 0
	is_left := true
	for left_len+s_len+right_len < length {
		if is_left {
			left_len += pad_len
			for left_len+s_len+right_len > length {
				left_len--
			}
		} else {
			right_len += pad_len
			for left_len+s_len+right_len > length {
				right_len--
			}
		}
		is_left = !is_left
	}

	left := LeftPad(s, padStr, left_len+s_len)
	return RightPad(left, padStr, left_len+s_len+right_len)
}

func Dedupe(s string) string {
	if len(s) <= 1 {
		return s
	}

	runes := []rune(s)
	result := make([]rune, 0, len(runes))

	// Add first character
	result = append(result, runes[0])

	// Compare each character with previous one
	for i := 1; i < len(runes); i++ {
		if runes[i] != runes[i-1] {
			result = append(result, runes[i])
		}
	}

	return string(result)
}

// ToWindowsPathSeparator converts a given file path to use Windows path separators.
// It handles empty paths, trims leading and trailing whitespace, and preserves network share paths.
//
// Parameters:
//   - path: The input file path as a string.
//
// Returns:
//   - A string representing the file path with Windows path separators.
//
// Example:
//
//	ToWindowsPathSeparator("C:/my_projects/s/s.go") // Returns "C:\my_projects\s\s.go"
func ToWindowsPathSeparator(path string) string {
	// Handle empty path
	if path == "" {
		return path
	}
	path = Trim(path)

	is_unc := false
	// Preserve network share path starting with \\
	if len(path) >= 2 && path[0] == '\\' && path[1] == '\\' {
		path = strings.TrimLeft(path[2:], "\\/")
		is_unc = true
	}

	// Replace all remaining contiguous separators with single backslash
	re := regexp.MustCompile(`[\\/]{2,}`)
	path = strings.ReplaceAll(re.ReplaceAllString(path, "\\"), "/", "\\")
	if is_unc {
		path = `\\` + path
	}
	return path
}

// ToLinuxPathSeparator converts a given file path to use Linux-style forward slashes as separators.
// It handles empty paths by returning them unchanged. Additionally, it replaces all contiguous
// separators (both backslashes and forward slashes) with a single forward slash.
//
// Parameters:
//   - path: The file path to be converted.
//
// Returns:
//   - A string representing the converted file path with Linux-style separators.
func ToLinuxPathSeparator(path string) string {
	// Handle empty path
	if path == "" {
		return path
	}

	// Replace all contiguous separators with single forward slash
	re := regexp.MustCompile(`[\\/]{2,}`)
	return strings.ReplaceAll(re.ReplaceAllString(path, "/"), "\\", "/")
}
