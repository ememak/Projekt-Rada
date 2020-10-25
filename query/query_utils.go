package query

import (
	"unicode"
)

func (t *PollSchema_QuestionType) IsValid() bool {
	return PollSchema_OPEN <= *t && *t <= PollSchema_CLOSE
}

func IsStringPrintable(s string) bool {
	for _, c := range s {
		if !unicode.IsGraphic(c) && !unicode.IsSpace(c) {
			return false
		}
	}
	return true
}
