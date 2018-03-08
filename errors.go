package anitogo

import "errors"

// Errors used throughout package
var (
	ErrorTokenNotFound                = errors.New("token not found in tokens.tokens slice")
	ErrorMatchingTokenNotFound        = errors.New("a matching token was not found")
	ErrorInvalidIndexDerived          = errors.New("the provided parameters resulted in an invalid index being derived, and thus a valid slice could not be returned")
	ErrorRuneNotFoundInString         = errors.New("the specified rune was not found in the specified string")
	ErrorTokenizationProducedNoTokens = errors.New("tokenization produced no tokens")
	ErrorInvalidEpisodeNumber = errors.New("the provided episode number string does not represent a valid episode number")
)
