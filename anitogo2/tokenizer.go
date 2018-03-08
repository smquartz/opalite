package anitogo2

import (
	"fmt"
	"regexp"

	"github.com/juju/errors"
)

func (t tokenizer) TokenizeByDelimiters(text string, enclosed bool) error {
	var delimiters string
	for _, delimiter := range t.Options.AllowedDelimiters {
		delimiters = delimiters + "\\" + string(delimiter)
	}

	pattern := fmt.Sprintf("([%v])", delimiters)
	cPattern, err := regexp.Compile(pattern)
	if err != nil {
		return errors.Trace(err)
	}
	splitText := cPattern.Split(text, -1)

	for _, subtext := range splitText {
		if subtext != "" {
			t.tokenManager.AddToken(tokenCategoryDelimiter, subtext, enclosed)
		} else {
			t.tokenManager.AddToken(tokenCategoryUnknown, subtext, enclosed)
		}
	}

	t.ValidateDelimiterTokens()

	return errors.Trace(nil)
}
