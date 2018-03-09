package anitogo

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/juju/errors"
)

func (t tokenizer) Tokenize() error {
	err := t.TokenizeByBrackets()
	if err != nil {
		return errors.Trace(err)
	}
	if t.tokenManager.Empty() {
		return errors.Trace(ErrorTokenizationProducedNoTokens)
	}
	return errors.Trace(nil)
}

func (t tokenizer) TokenizeByBrackets() error {
	brackets := [][]rune{
		[]rune{'(', ')'},           // U+0028-U+0029 Parenthesis
		[]rune{'[', ']'},           //  U+005B-U+005D Square bracket
		[]rune{'{', '}'},           // U+007B-U+007D Curly bracket
		[]rune{'\u300C', '\u300D'}, // Corner bracket
		[]rune{'\u300E', '\u300F'}, // White corner bracket
		[]rune{'\u3010', '\u3011'}, // Black lenticular bracket
		[]rune{'\uFF08', '\uFF09'}, // Fullwidth parenthesis
	}

	var openBrackets []rune
	for _, bracketPair := range brackets {
		openBrackets = append(openBrackets, bracketPair[0])
	}

	text := t.Filename
	isBracketOpen := false
	var matchingBracket rune

	findFirstBracket := func() (int, rune) {
		var index int
		for idx, bracket := range text {
			var found bool
			for _, r := range openBrackets {
				if bracket == r {
					found = true
				}
			}
			if found {
				index = idx
				break
			}
			index = -1
		}

		var matchingBracket rune
		for _, bracketPair := range brackets {
			if bracketPair[0] == rune(text[index]) {
				matchingBracket = bracketPair[1]
				break
			}
		}

		return index, matchingBracket
	}

	for idx, r := range text {
		var bracketIndex int
		var matchingBracket rune

		if !isBracketOpen {
			bracketIndex, matchingBracket = findFirstBracket()
		} else {
			// Looking for the matching bracket allows us to better handle
			// some rare cases with nested brackets.
			bracketIndex = strings.IndexRune(text, matchingBracket)
		}

		// found a token before the bracket... maybe?
		if bracketIndex != 0 {
			if bracketIndex > -1 {
				t.TokenizeByPreidentified(text[:bracketIndex], isBracketOpen)
			} else {
				t.TokenizeByPreidentified(text, isBracketOpen)
			}
		}

		// found bracket
		if bracketIndex > -1 {
			t.tokenManager.AddToken(tokenCategoryBracket, string(text[bracketIndex]), true)

			isBracketOpen = !isBracketOpen
			text = text[bracketIndex+1:]
		} else {
			log.Println("Clearing text where len = " + strconv.Itoa(len(text)) + " in TokenizeByBrackets")
			text = "" // consider this
		}
	}

	return errors.Trace(nil)
}

func (t tokenizer) TokenizeByPreidentified(text string, enclosed bool) error {
	preidentifiedTokens, err := t.keywordManager.Peek(text)
	if err != nil {
		return errors.Trace(err)
	}

	lastTokenEndIndex := 0
	for _, index := range preidentifiedTokens {
		if lastTokenEndIndex != index.Begin {
			// tokenize text between preidentified tokens
			err = t.TokenizeByDelimiters(text[lastTokenEndIndex:index.Begin], enclosed)
			if err != nil {
				return errors.Trace(err)
				// continue
			}

			t.tokenManager.AddToken(tokenCategoryIdentifier, text[index.Begin:index.End], enclosed)
			lastTokenEndIndex = index.End
		}
	}

	if lastTokenEndIndex != len(text) {
		// tokenize text after preidentified tokens or all text if no preidentified
		// tokens
		err = t.TokenizeByDelimiters(text[lastTokenEndIndex:], enclosed)
		if err != nil {
			return errors.Trace(err)
		}
	}

	return errors.Trace(nil)
}

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
