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

		for _, opt := range qa.Options {
			sum := 0
			if !IsStringPrintable(opt.Name) {
				return fmt.Errorf("Error! Answer option contains non valid characters.")
			}
			if opt.Selected {
				sum += 1
			}
			if qa.Type == PollSchema_CLOSE {
				if sum == 0 {
					return fmt.Errorf("Error! Answer option is not selected.")
				}
				if sum > 1 {
					return fmt.Errorf("Error! Multiple answer options are selected.")
				}
			}
		}

		if !qa.Type.IsValid() {
			return fmt.Errorf("Error! Wrong question type.")
		}

		if !IsStringPrintable(qa.Answer) {
			return fmt.Errorf("Error! Answer contains non valid characters.")
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
