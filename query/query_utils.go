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
			return fmt.Errorf("Error! Question contains non valid characters.")
		}

		if !qa.Type.IsValid() {
			return fmt.Errorf("Error! Wrong question type.")
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
