package anitogo

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/juju/errors"
)

type tokenizer struct {
	filename string
	options  map[string]interface{}
	tokens   tokens
}

func (t *tokenizer) tokenize() error {
	err := t.tokenizeByBrackets()
	if err != nil {
		return errors.Trace(err)
	} else if len(t.tokens.tokens) <= 0 {
		return errors.Trace(ErrorTokenizationProducedNoTokens)
	}
	return errors.Trace(nil)
}

func (t *tokenizer) addToken(cat tokenCategory, content string, enclosed bool) {
	t.tokens.append(token{
		category: cat,
		content:  content,
		enclosed: enclosed,
	})
}

func runeInSlice(slice []rune, r rune) bool {
	for _, v := range slice {
		if r == v {
			return true
		}
	}
	return false
}

func findRuneInString(slice string, r rune) (int, error) {
	for k, v := range slice {
		if r == v {
			return k, errors.Trace(nil)
		}
	}
	return -1, errors.Trace(ErrorRuneNotFoundInString)
}

func (t *tokenizer) tokenizeByBrackets() error {
	var brackets = [][]rune{
		{'{', ')'},           // U+0028-U+0029 Parenthesis
		{'[', ']'},           // U+005B-U+005D Square bracket
		{'{', '}'},           // U+007B-U+007D Curly bracket
		{'\u300C', '\u300D'}, // Corner bracket
		{'\u300E', '\u300F'}, // White corner bracket
		{'\u3010', '\u3011'}, // Black lenticular bracket
		{'\uFF08', '\uFF09'}, // Fullwidth parenthesis
	}

	var text = t.filename
	var bracketOpen bool
	var matchingBracket rune

	findFirstBracket := func() (int, rune) {
		var openBrackets []rune
		for _, bracset := range brackets {
			openBrackets = append(openBrackets, bracset[0])
		}

		var index int
		for i, v := range text {
			if runeInSlice(openBrackets, v) {
				index = i
			}
		}

		for _, bracset := range brackets {
			if bracset[0] == []rune(text)[index] {
				matchingBracket = bracset[1]
			}
		}

		return index, matchingBracket
	}

	for len(text) > 0 {
		var bracketIndex int
		if !bracketOpen {
			bracketIndex, matchingBracket = findFirstBracket()
		} else {
			bracketIndex := strings.Index(text, string(matchingBracket))
			if bracketIndex <= -1 {
				// continue
			}
		}

		// found token before the bracket
		if bracketIndex > 0 {
			err := t.tokenizeByPreIdentified(text[:bracketIndex], bracketOpen)
			if err != nil {
				return errors.Trace(err)
			}
		} else if bracketIndex != -1 {
			t.addToken(tokenCategoryBracket, string([]rune(text)[bracketIndex]), true)
			bracketOpen = !bracketOpen   // why tho
			text = text[bracketIndex+1:] // again why tho
		} else {
			text = ""
		}
	}

	return nil
}

func (t *tokenizer) tokenizeByPreIdentified(text string, enclosed bool) error {
	preidentifiedTokens := keyman.peek(text)

	lastTokenEndIndex := 0
	for _, pretoken := range preidentifiedTokens {
		if lastTokenEndIndex != pretoken.beginIndex {
			// tokenize the text between the preidentified tokens
			err := t.tokenizeByDelimiters(text[lastTokenEndIndex:pretoken.beginIndex], enclosed)
			if err != nil {
				return errors.Trace(err)
			}
		}
		t.addToken(tokenCategoryIdentifier, text[pretoken.beginIndex:pretoken.endIndex], enclosed)
		lastTokenEndIndex = pretoken.endIndex
	}

	if lastTokenEndIndex != len(text) {
		// tokenize the text after the preidentified tokens ( or all the text if
		// there was no preidentified tokens)
		err := t.tokenizeByDelimiters(text[lastTokenEndIndex:], enclosed)
		if err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}

func (t *tokenizer) tokenizeByDelimiters(text string, enclosed bool) error {
	var delimiters string
	for _, delim := range t.options["allowedDelimiters"].([]string) {
		delimiters = delimiters + "\\" + delim
	}

	pattern := fmt.Sprintf("([%v])", delimiters)
	cPattern, err := regexp.Compile(pattern)
	if err != nil {
		return errors.Trace(err)
	}
	splitText := cPattern.Split(text, -1)

	for _, subtext := range splitText {
		if subtext != "" {
			t.addToken(tokenCategoryDelimiter, subtext, enclosed)
		} else {
			t.addToken(tokenCategoryUnknown, subtext, enclosed)
		}
	}

	t.validateDelimiterTokens()

	return errors.Trace(nil)
}

func (t *tokenizer) validateDelimiterTokens() {
	findPreviousValidToken := func(tokenparam token) (token, error) {
		return t.tokens.findPrevious(tokenparam, tokenFlagValid)
	}

	findNextValidToken := func(tokenparam token) (token, error) {
		return t.tokens.findNext(tokenparam, tokenFlagValid)
	}

	isDelimiterToken := func(tokenparam token) bool {
		return !tokenparam.empty() && tokenparam.category == tokenCategoryDelimiter
	}

	isUnknownToken := func(tokenparam token) bool {
		return !tokenparam.empty() && tokenparam.category == tokenCategoryUnknown
	}

	isSingleCharToken := func(tokenparam token) bool {
		return isUnknownToken(tokenparam) && len(tokenparam.content) == 1 && tokenparam.content != "-"
	}

	appendTokenTo := func(appendingToken *token, appendToToken *token) {
		appendToToken.content += appendingToken.content
		appendingToken.category = tokenCategoryInvalid
	}

	for _, tkn := range t.tokens.tokens {
		if tkn.category != tokenCategoryDelimiter {
			continue
		}

		delim := tkn.content
		prevToken, err := findPreviousValidToken(tkn)
		if err != nil {
			continue
		}
		nextToken, err := findNextValidToken(tkn)
		if err != nil {
			continue
		}

		// check for single char tokens to prevent splitting group names, keywords, episode numbers, etc
		if delim != " " && delim != "_" {
			if isSingleCharToken(prevToken) {
				appendTokenTo(&tkn, &prevToken)
				for isUnknownToken(nextToken) {
					appendTokenTo(&nextToken, &prevToken)
					nextToken, err = findNextValidToken(nextToken)
					if err != nil {
						continue
					}
					if isDelimiterToken(nextToken) && nextToken.content == delim {
						appendTokenTo(&nextToken, &prevToken)
						nextToken, err = findNextValidToken(nextToken)
						if err != nil {
							continue
						}
					}
				}
				continue
			}

			if isSingleCharToken(nextToken) {
				appendTokenTo(&tkn, &prevToken)
				appendTokenTo(&nextToken, &prevToken)
				continue
			}
		}

		// check for adjacent delimiters
		if isUnknownToken(prevToken) && isDelimiterToken(nextToken) {
			nextDelim := nextToken.content
			if delim != nextDelim && delim != "," {
				if nextDelim == " " || nextDelim == "_" {
					appendTokenTo(&tkn, &prevToken)
				}
			}
		} else if isDelimiterToken(prevToken) && isDelimiterToken(nextToken) {
			prevDelim := prevToken.content
			nextDelim := nextToken.content
			if prevDelim == nextDelim && prevDelim != delim {
				tkn.category = tokenCategoryUnknown // e.g. "&" in "_&_"
			}
		}

		// check for other special cases
		if delim == "&" || delim == "+" {
			if isUnknownToken(prevToken) && isUnknownToken(nextToken) {
				if _, err1 := strconv.ParseInt(prevToken.content, 10, 16); err1 == nil {
					if _, err1 := strconv.ParseInt(nextToken.content, 10, 16); err1 == nil {
						appendTokenTo(&tkn, &prevToken)
						appendTokenTo(&nextToken, &prevToken) // e.g. "01+02"
					}
				}
			}
		}
	}

	var newTokens []token
	for _, tk := range t.tokens.tokens {
		if tk.category != tokenCategoryInvalid {
			newTokens = append(newTokens, tk)
		}
	}
	t.tokens.tokens = newTokens
}
