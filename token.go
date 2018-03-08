package anitogo

import (
	"github.com/juju/errors"
)

type tokenCategory int

const (
	tokenCategoryUnknown tokenCategory = iota
	tokenCategoryBracket
	tokenCategoryDelimiter
	tokenCategoryIdentifier
	tokenCategoryInvalid
)

type tokenFlag int

const (
	tokenFlagNone    tokenFlag = 0
	tokenFlagBracket tokenFlag = 1 << iota
	tokenFlagNotBracket
	tokenFlagDelimiter
	tokenFlagNotDelimiter
	tokenFlagIdentifier
	tokenFlagNotIdentifier
	tokenFlagUnknown
	tokenFlagNotUnknown
	tokenFlagValid
	tokenFlagNotValid
	tokenFlagEnclosed
	tokenFlagNotEnclosed

	tokenFlagMaskCategories = tokenFlagBracket | tokenFlagNotBracket | tokenFlagDelimiter | tokenFlagNotDelimiter | tokenFlagIdentifier | tokenFlagNotIdentifier | tokenFlagUnknown | tokenFlagNotUnknown | tokenFlagValid | tokenFlagNotValid
	tokenFlagMaskEnclosed   = tokenFlagEnclosed | tokenFlagNotEnclosed
)

type token struct {
	category tokenCategory
	content  string
	enclosed bool
}

func newToken(category tokenCategory, content string, enclosed bool) *token {
	return &token{
		category: category,
		content:  content,
		enclosed: enclosed,
	}
}

// checkFlag returns whether flags contains flag
func checkFlag(flags, flag tokenFlag) bool {
	return (flags & flag) == flag
}

// checkCategory returns categories contains category
func checkCategory(categories, category tokenCategory) bool {
	return (categories & category) == category
}

// checks whether token is null token
func (t token) empty() bool {
	var emptyToken token
	return t == emptyToken
}

// checkFlags does what exactly?
func (t token) checkFlags(flags tokenFlag) bool {
	// not sure what this does exactly
	checkCategory := func(fe, fn tokenFlag, cat tokenCategory) bool {
		if checkFlag(flags, fe) {
			return t.category == cat
		} else if checkFlag(flags, fn) {
			return t.category != cat
		}
		return false
	}

	// if contains tokenFlagMaskEnclosed
	if checkFlag(flags, tokenFlagMaskEnclosed) {
		var success bool
		if checkFlag(flags, tokenFlagEnclosed) {
			success = t.enclosed
		} else {
			success = !t.enclosed
		}
		if !success {
			return false
		}
	}

	if checkFlag(flags, tokenFlagMaskCategories) {
		if checkCategory(tokenFlagBracket, tokenFlagNotBracket, tokenCategoryBracket) {
			return true
		}
		if checkCategory(tokenFlagDelimiter, tokenFlagNotDelimiter, tokenCategoryDelimiter) {
			return true
		}
		if checkCategory(tokenFlagIdentifier, tokenFlagNotIdentifier, tokenCategoryIdentifier) {
			return true
		}
		if checkCategory(tokenFlagUnknown, tokenFlagNotUnknown, tokenCategoryUnknown) {
			return true
		}
		if checkCategory(tokenFlagNotValid, tokenFlagValid, tokenCategoryInvalid) {
			return true
		}
		return false
	}

	return true
}

type tokens struct {
	// instance tokens // what does instance field do even
	tokens []token
}

func (t tokens) empty() bool {
	return len(t.tokens) == 0
}

func (t *tokens) append(token token) {
	t.tokens = append(t.tokens, token)
}

func (t tokens) getList(flags tokenFlag, beginToken, endToken token) ([]token, error) {
	tokens := t.tokens
	var emptyToken token
	var err error
	var beginIndex int
	var endIndex int
	if beginToken == emptyToken {
		beginIndex = 0
	} else {
		beginIndex, err = t.getIndex(beginToken)
		if err != nil {
			return []token{}, errors.Trace(err)
		}
	}
	if endToken == emptyToken {
		endIndex = len(tokens)
	} else {
		endIndex, err = t.getIndex(endToken)
		if err != nil {
			return []token{}, errors.Trace(err)
		}
	}

	// catch potential out of range index case
	if (endIndex + 1) > (len(tokens) - 1) {
		return []token{}, errors.Trace(ErrorInvalidIndexDerived)
	}
	tokens = tokens[beginIndex : endIndex+1]
	if flags == 0 {
		return tokens, nil
	}
	var newTokens []token
	for _, value := range tokens {
		if value.checkFlags(flags) {
			newTokens = append(newTokens, value)
		}
	}
	return newTokens, errors.Trace(nil)
}

func (t tokens) getIndex(token token) (int, error) {
	for index, value := range t.tokens {
		if token == value {
			return index, errors.Trace(nil)
		}
	}
	return 0, errors.Trace(ErrorTokenNotFound)
}

// distance returns the index distance between two tokens in a tokens{}'s tokens field
func (t tokens) distance(tokenStart, tokenEnd token) (int, error) {
	var emptyToken token
	var startIndex int
	var err error
	if tokenStart == emptyToken {
		startIndex = 0
	} else {
		startIndex, err = t.getIndex(tokenStart)
		if err != nil {
			return 0, errors.Trace(err)
		}
	}
	var endIndex int
	if tokenEnd == emptyToken {
		endIndex = len(t.tokens)
	} else {
		endIndex, err = t.getIndex(tokenEnd)
		if err != nil {
			return 0, errors.Trace(err)
		}
	}
	return endIndex - startIndex, errors.Trace(nil)
}

func findInTokens(tokens []token, flags tokenFlag) (token, error) {
	for _, token := range tokens {
		if token.checkFlags(flags) {
			return token, errors.Trace(nil)
		}
	}
	return token{}, errors.Trace(ErrorMatchingTokenNotFound)
}

// findInTokens finds a token in a tokens{}'s tokens field, that satisfies token.checkFlags(flags) where flags is a parameter
func (t tokens) find(flags tokenFlag) (token, error) {
	return findInTokens(t.tokens, flags)
}

func reverseTokensSlice(tokens []token) []token {
	lastIndex := len(tokens) - 1
	for i := 0; i < len(tokens)/2; i++ {
		tokens[i], tokens[lastIndex-i] = tokens[lastIndex-i], tokens[i]
	}
	return tokens
}

func (t tokens) findPrevious(tokenParam token, flags tokenFlag) (token, error) {
	var emptyToken token
	var tokens = t.tokens
	if tokenParam == emptyToken {
		tokens = reverseTokensSlice(tokens)
	} else {
		tokenIndex, err := t.getIndex(tokenParam)
		if err != nil {
			return token{}, errors.Trace(err)
		}
		tokens = reverseTokensSlice(tokens[tokenIndex-1:])
	}
	return findInTokens(tokens, flags)
}

func (t tokens) findNext(tokenParam token, flags tokenFlag) (token, error) {
	tokens := t.tokens
	var emptyToken token
	if tokenParam != emptyToken {
		tokenIndex, err := t.getIndex(tokenParam)
		if err != nil {
			return emptyToken, errors.Trace(err)
		}
		tokens = tokens[tokenIndex+1:]
	}
	return findInTokens(tokens, flags)
}
