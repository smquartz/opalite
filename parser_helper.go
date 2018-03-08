package anitogo

import (
	"log"
	"strconv"
	"strings"
	"unicode"

	"github.com/juju/errors"
)

const dashes = "-\u2010\u2011\u2012\u2013\u2014\u2015"

func findNumberInString(str string) int {
	for _, c := range str {
		if unicode.IsDigit(c) {
			return strings.IndexRune(str, c)
		}
	}
	return -1
}

func findNonNumberInString(str string) int {
	for _, r := range str {
		if !unicode.IsDigit(r) {
			return strings.IndexRune(str, r)
		}
	}
	return -1
}

func isHexadecimalString(str string) bool {
	_, err := strconv.ParseInt(str, 16, 64)
	return err == nil
}

func getNumberFromOrdinal(str string) (uint, error) {
	ordinals := map[string]uint{
		"1st": 1, "first": 1,
		"2nd": 2, "second": 2,
		"3rd": 3, "third": 3,
		"4th": 4, "fourth": 4,
		"5th": 5, "fifth": 5,
		"6th": 6, "sixth": 6,
		"7th": 7, "seventh": 7,
		"8th": 8, "eigth": 8,
		"9th": 9, "ninth": 9,
	}

	lowerStr := strings.ToLower(str)
	num, ok := ordinals[lowerStr]
	if !ok {
		return 0, errors.Trace(ErrorUnrecognisedOrdinal)
	}
	return num, errors.Trace(nil)
}

func isCRC32(str string) bool {
	return len(str) == 8 && isHexadecimalString(str)
}

func isDashCharacter(str string) bool {
	if len(str) != 1 {
		return false
	}
	for _, dash := range dashes {
		if str == string(dash) {
			return true
		}
	}
	return false
}

func isLatinRune(r rune) bool {
	return unicode.In(r, unicode.Latin)
}

func isMostlyLatinString(str string) bool {
	if len(str) <= 0 {
		return false
	}
	latinLength := 0
	for _, r := range str {
		if isLatinRune(r) {
			latinLength++
		}
	}
	return (float32(latinLength) / float32(len(str))) >= 0.5
}

func (t *tokens) checkAnimeSeasonKeyword(tkn token) (bool, error) {
	setAnimeSeason := func(first, second token, content string) {
		elems.insert(elementCategoryAnimeSeason, content)
		first.category = tokenCategoryIdentifier
		second.category = tokenCategoryIdentifier
	}

	prevToken, err := t.findPrevious(tkn, tokenFlagNotDelimiter)
	if err != nil {
		return false, errors.Trace(err)
	}
	if !prevToken.empty() {
		num, err := getNumberFromOrdinal(prevToken.content)
		if err == nil {
			setAnimeSeason(prevToken, tkn, strconv.FormatUint(uint64(num), 10))
			return true, errors.Trace(nil)
		}
		log.Println("Uh oh " + err.Error())
	}

	nextToken, err := t.findNext(tkn, tokenFlagNotDelimiter)
	if err != nil {
		return false, errors.Trace(err)
	}
	if _, err1 := strconv.ParseInt(nextToken.content, 10, 64); !nextToken.empty() && err1 == nil {
		setAnimeSeason(tkn, nextToken, nextToken.content)
		return true, errors.Trace(nil)
	}
	return false, errors.Trace(nil)
}

func (t tokens) isTokenIsolated(tkn token) (bool, error) {
	prevToken, err := t.findPrevious(tkn, tokenFlagNotDelimiter)
	if err != nil {
		return false, errors.Trace(ErrorMatchingTokenNotFound)
	}
	if prevToken.category != tokenCategoryBracket {
		return false, errors.Trace(nil)
	}

	nextToken, err := t.findNext(tkn, tokenFlagNotDelimiter)
	if err != nil {
		return false, errors.Trace(ErrorMatchingTokenNotFound)
	}
	if nextToken.category != tokenCategoryBracket {
		return false, errors.Trace(nil)
	}

	return true, errors.Trace(nil)
}

func (t tokens) buildElement(cat elementCategory, tokenBegin, tokenEnd token, keepDelimiters bool) error {
	element := ""

	tkns, err := t.getList(tokenFlagNone, tokenBegin, tokenEnd)
	if err != nil {
		return errors.Trace(err)
	}

	for _, tkn := range tkns {
		if tkn.category == tokenCategoryUnknown {
			element += tkn.content
		} else if tkn.category == tokenCategoryBracket {
			element += tkn.content
		} else if tkn.category == tokenCategoryDelimiter {
			delimiter := tkn.content
			if keepDelimiters {
				element += delimiter
			} else if tkn != tokenBegin && tkn != tokenEnd {
				if delimiter == "," || delimiter == "&" {
					element += delimiter
				} else {
					element += " "
				}
			}
		}
	}

	if !keepDelimiters {
		for _, delimrune := range dashes + " " {
			delim := string(delimrune)
			element = strings.Replace(element, delim, "", -1)
		}
	}

	if element != "" {
		elems.insert(cat, element)
	}

	return errors.Trace(nil)
}
