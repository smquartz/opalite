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

func numberComesAfterPrefix(cat category, tkn *token) bool {
	numberBegin := findNumberInString(tkn.Content)
	prefix := tkn.Content[:numberBegin]
}

func (p *parser) SetAnimeSeason(first, second *token, content string) error {
	contentNum, err := strconv.ParseUint(content, 10, 16)
	if err != nil {
		return errors.Trace(err)
	}
	p.AnimeFile.AnimeSeason = append(p.AnimeFile.AnimeSeason, uint(contentNum))
	first.Category = tokenCategoryIdentifier
	second.Category = tokenCategoryIdentifier
	return errors.Trace(nil)
}

func (p *parser) CheckAnimeSeasonKeyword(tkn *token) (bool, error) {
	prevToken, err := p.tokenizer.tokenManager.tokens.FindPrevious(*tkn, tokenFlagsNotDelimiter)
	if err != nil {
		return false, errors.Trace(err)
	}
	if !prevToken.Empty() {
		num, err := getNumberFromOrdinal(prevToken.Content)
		if err == nil {
			p.SetAnimeSeason(prevToken, tkn, strconv.FormatUint(uint64(num), 10))
			return true, errors.Trace(nil)
		}
		log.Println("Uh oh " + err.Error())
	}

	nextToken, err := p.tokenizer.tokenManager.tokens.FindNext(*tkn, tokenFlagsNotDelimiter)
	if err != nil {
		return false, errors.Trace(err)
	}
	if _, err1 := strconv.ParseInt(nextToken.Content, 10, 64); !nextToken.Empty() && err1 == nil {
		p.SetAnimeSeason(tkn, nextToken, nextToken.Content)
		return true, errors.Trace(nil)
	}
	return false, errors.Trace(nil)
}

func (t tokens) IsTokenIsolated(tkn token) (bool, error) {
	prevToken, err := t.FindPrevious(tkn, tokenFlagsNotDelimiter)
	if err != nil {
		return false, errors.Trace(ErrorMatchingTokenNotFound)
	}
	if prevToken.Category != tokenCategoryBracket {
		return false, errors.Trace(nil)
	}

	nextToken, err := t.FindNext(tkn, tokenFlagsNotDelimiter)
	if err != nil {
		return false, errors.Trace(ErrorMatchingTokenNotFound)
	}
	if nextToken.Category != tokenCategoryBracket {
		return false, errors.Trace(nil)
	}

	return true, errors.Trace(nil)
}

func (p *parser) BuildElement(cat category, tokenBegin, tokenEnd token, keepDelimiters bool) error {
	element := ""

	tkns, err := p.tokenizer.tokenManager.tokens.getList(tokenFlagsNone, tokenBegin, tokenEnd)
	if err != nil {
		return errors.Trace(err)
	}

	for _, tkn := range tkns {
		if tkn.Category == tokenCategoryUnknown {
			element += tkn.Content
		} else if tkn.Category == tokenCategoryBracket {
			element += tkn.Content
		} else if tkn.Category == tokenCategoryDelimiter {
			delimiter := tkn.Content
			if keepDelimiters {
				element += delimiter
			} else if *tkn != tokenBegin && *tkn != tokenEnd {
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

	switch cat {
	case categoryAnimeSeason:
		num, err := strconv.ParseUint(element, 10, 16)
		if err != nil {
			return errors.Trace(err)
		}
		p.AnimeFile.AnimeSeason = append(p.AnimeFile.AnimeSeason, uint(num))
	}

	return errors.Trace(nil)
}
