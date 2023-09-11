package sqlbuilder

import "strings"

const (
	cacheQuestionMarkN = 100
)

var (
	questionMarksCache = make(map[int]string, cacheQuestionMarkN)
)

func buildQuestionMarks(n int) string {
	if n <= 0 {
		return "()"
	}
	b := &strings.Builder{}
	b.Grow(2*n + 1)
	b.WriteByte(openParentheses)
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(comma)
		}
		b.WriteByte(questionMark)
	}
	b.WriteByte(closeParentheses)

	return b.String()
}

func QuestionMarks(n int) string {
	if 1 <= n && n <= cacheQuestionMarkN {
		return questionMarksCache[n]
	}
	return buildQuestionMarks(n)
}

func init() {
	for i := 1; i <= cacheQuestionMarkN; i++ {
		questionMarksCache[i] = buildQuestionMarks(i)
	}
}
