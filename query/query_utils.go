package query

import (
	"fmt"
	"unicode"
)

func (t *PollSchema_QuestionType) IsValid() bool {
	return PollSchema_OPEN <= *t && *t <= PollSchema_CLOSE
}

func (t *PollSchema) IsValid() error {
	for _, qa := range t.Questions {
		if !IsStringPrintable(qa.Question) {
			return fmt.Errorf("Error! Question contains invalid characters.")
		}

		for _, opt := range qa.Options {
			if !IsStringPrintable(opt) {
				return fmt.Errorf("Error! Answer option contains invalid characters.")
			}
		}

		if !qa.Type.IsValid() {
			return fmt.Errorf("Error! Wrong question type.")
		}

		for _, ans := range qa.Answers {
			if !IsStringPrintable(ans) {
				return fmt.Errorf("Error! Answer contains invalid characters.")
			}
		}
	}
	return nil
}

func IsStringPrintable(s string) bool {
	for _, c := range s {
		if !unicode.IsGraphic(c) && !unicode.IsSpace(c) {
			return false
		}
	}
	return true
}
