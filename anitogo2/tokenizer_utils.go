package anitogo2

import "strconv"

func newTokenizer(filename string, options TokenizerOptions, tkman *tokenManager, keyman *keywordManager) *tokenizer {
	tknzr := new(tokenizer)
	tknzr.Filename = filename
	tknzr.Options = options
	tknzr.tokenManager = tkman
	tknzr.keywordManager = keyman
	return tknzr
}

func newTokenizerWithManagers(filename string, options TokenizerOptions) *tokenizer {
	keyman := newKeywordManager()
	tkman := newTokenManager()
	return newTokenizer(filename, options, tkman, keyman)
}

func newTokenizerWithDefaults(filename string) *tokenizer {
	return newTokenizerWithManagers(filename, TokenizerDefaultOptions)
}

func (t tokenizer) FindPreviousValidToken(tkn token) (*token, error) {
	return t.tokenManager.tokens.FindPrevious(tkn, tokenFlagsValid)
}

func (t tokenizer) FindNextValidToken(tkn token) (*token, error) {
	return t.tokenManager.tokens.FindNext(tkn, tokenFlagsValid)
}

func isDelimiterToken(tkn token) bool {
	return !tkn.Empty() && tkn.Category == tokenCategoryDelimiter
}

func isUnknownToken(tkn token) bool {
	return !tkn.Empty() && tkn.Category == tokenCategoryUnknown
}

func isSingleCharacterToken(tkn token) bool {
	return isUnknownToken(tkn) && len(tkn.Content) == 1 && tkn.Content != "-"
}

func (t *token) AppendToken(tkn *token) {
	t.Content += tkn.Content
	tkn.Category = tokenCategoryInvalid
}

func (t *tokenizer) ValidateDelimiterTokens() {
	for _, tkn := range t.tokenManager.tokens {
		if tkn.Category != tokenCategoryDelimiter {
			continue
		}

		delimiter := tkn.Content
		prevToken, err := t.FindPreviousValidToken(*tkn)
		if err != nil {
			continue
			// is this actually best
		}
		nextToken, err := t.FindNextValidToken(*tkn)
		if err != nil {
			continue
			// is this actually best
		}

		// check for single character tokens to prevent splitting up: group names, keywords, episode numbers, etc
		if delimiter != " " && delimiter != "_" {
			if isSingleCharacterToken(*prevToken) {
				prevToken.AppendToken(tkn)

				for isUnknownToken(*nextToken) {
					prevToken.AppendToken(nextToken)
					nextToken, err := t.FindNextValidToken(*nextToken)
					if err != nil {
						continue
						// is this the best thing to do
					}

					if isDelimiterToken(*nextToken) && nextToken.Content == delimiter {
						prevToken.AppendToken(nextToken)
						nextToken, err = t.FindNextValidToken(*nextToken)
						if err != nil {
							continue
							// is this the best thing to do
						}
					}
				}
				continue
			}

			if isSingleCharacterToken(*nextToken) {
				prevToken.AppendToken(tkn)
				prevToken.AppendToken(nextToken)
				continue
			}
		}

		// check for adjacent delimiters
		if isUnknownToken(*prevToken) && isDelimiterToken(*nextToken) {
			nextDelimiter := nextToken.Content
			if delimiter != nextDelimiter && delimiter != "," {
				if nextDelimiter == " " || nextDelimiter == "_" {
					prevToken.AppendToken(tkn)
				}
			} else if isDelimiterToken(*prevToken) && isDelimiterToken(*nextToken) {
				prevDelimiter := prevToken.Content
				nextDelimiter = nextToken.Content
				if prevDelimiter == nextDelimiter && prevDelimiter != delimiter {
					tkn.Category = tokenCategoryUnknown // e.g. "&" in "_&_"
				}
			}

			// check for other special cases
			if delimiter == "&" || delimiter == "+" {
				if isUnknownToken(*prevToken) && isUnknownToken(*nextToken) {
					if _, err := strconv.ParseInt(prevToken.Content, 10, 16); err == nil {
						if _, err := strconv.ParseInt(nextToken.Content, 10, 16); err == nil {
							prevToken.AppendToken(tkn)
							prevToken.AppendToken(nextToken)
						}
					}
				}
			}
		}
	}

	var newTokens tokens
	for _, tk := range t.tokenManager.tokens {
		if tk.Category != tokenCategoryInvalid {
			newTokens = append(newTokens, tk)
		}
	}
	t.tokenManager.tokens = newTokens
}
