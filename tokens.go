package anitogo

import "github.com/juju/errors"

func (t tokens) getList(flags tokenFlags, beginToken, endToken token) (tokens, error) {
	// var declaration
	var err error
	var beginIndex int
	var endIndex int

	// if no beginToken provided, start at the start
	if beginToken.Empty() {
		beginIndex = 0
	} else {
		// else start at beginToxen's index
		beginIndex, err = t.GetIndex(beginToken)
		if err != nil {
			return tokens{}, errors.Trace(err)
		}
	}

	// if no endToken provided, finish at the end
	if endToken.Empty() {
		endIndex = len(t)
	} else {
		// else, finish at endToken's index
		endIndex, err = t.GetIndex(endToken)
		if err != nil {
			return tokens{}, errors.Trace(err)
		}
	}

	t = t[beginIndex : endIndex+1]
	if flags == tokenFlagsNone {
		return t, errors.Trace(nil)
	}
	var newt tokens
	for _, tkn := range t {
		if tkn.HasFlags(flags) {
			newt = append(newt, tkn)
		}
	}
	return newt, errors.Trace(nil)
}

func (tm tokenManager) AddToken(cat tokenCategory, content string, enclosed bool) {
	tkn := &token{Category: cat, Content: content, Enclosed: enclosed}
	tm.tokens = append(tm.tokens, tkn)
}
