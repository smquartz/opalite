package anitogo

import "github.com/juju/errors"

// newTokenManager returns a new tokenManager instance
func newTokenManager() *tokenManager {
	return new(tokenManager)
}

// Empty returns whether the token is in its default state
func (t token) Empty() bool {
	return t == token{}
}

// GetIndex returns the index of a token
func (t tokens) GetIndex(tkn token) (int, error) {
	for i, tk := range t {
		if *tk == tkn {
			return i, errors.Trace(nil)
		}
	}
	return -1, errors.Trace(ErrorTokenNotFound)
}

// Distance returns the distance between the indexes of two tokens
func (t tokens) Distance(tokenStart, tokenEnd token) (int, error) {
	var startIndex int
	var endIndex int
	var err error

	// if no start token provided, start at beginning
	if tokenStart.Empty() {
		startIndex = 0
	} else {
		// else, starting at tokenStart's index
		startIndex, err = t.GetIndex(tokenStart)
		if err != nil {
			return -1, errors.Trace(err)
		}
	}

	// if no end token provided, finish at end
	if tokenEnd.Empty() {
		endIndex = len(t)
	} else {
		// else, finish at tokenEnd's index
		endIndex, err = t.GetIndex(tokenEnd)
		if err != nil {
			return -1, errors.Trace(err)
		}
	}

	return endIndex - startIndex, errors.Trace(nil)
}

// FindInTokens finds a token in a tokens slice containing the specified flags
func (t tokens) FindInTokens(flags tokenFlags) (*token, error) {
	for _, tkn := range t {
		if tkn.HasFlags(flags) {
			return tkn, errors.Trace(nil)
		}
	}
	return &token{}, errors.Trace(ErrorMatchingTokenNotFound)
}

// ReverseOrder reverses the order of a tokens slice
func (t *tokens) ReverseOrder() {
	lastIndex := len(*t) - 1
	var tkns = *t
	for i := 0; i < len(*t)/2; i++ {
		tkns[i], tkns[lastIndex-i] = tkns[lastIndex-i], tkns[i]
	}
	t = &tkns
}

func (t tokens) FindPrevious(tkn token, flags tokenFlags) (*token, error) {
	if tkn.Empty() {
		t.ReverseOrder()
	} else {
		tokenIndex, err := t.GetIndex(tkn)
		if err != nil {
			return &token{}, errors.Trace(err)
		}
		t = t[tokenIndex-1:]
		t.ReverseOrder()
	}
	return t.FindInTokens(flags)
}

func (t tokens) FindNext(tkn token, flags tokenFlags) (*token, error) {
	if tkn.Empty() {
		tokenIndex, err := t.GetIndex(tkn)
		if err != nil {
			return &token{}, errors.Trace(err)
		}
		t = t[tokenIndex+1:]
		t.ReverseOrder()
	}
	return t.FindInTokens(flags)
}
