package anitogo2

// TokenizerOptions configures the tokenizer's behaviour
type TokenizerOptions struct {
	AllowedDelimiters  string
	IgnoredStrings     []string
	ParseEpisodeNumber bool
	ParseEpisodeTitle  bool
	ParseFileExtension bool
	ParseReleaseGroup  bool
}

// TokenizerDefaultOptions provides sane default parser configuration options
var TokenizerDefaultOptions = TokenizerOptions{
	AllowedDelimiters:  "_.&+,|",
	IgnoredStrings:     []string{},
	ParseEpisodeNumber: true,
	ParseEpisodeTitle:  true,
	ParseFileExtension: true,
	ParseReleaseGroup:  true,
}

type tokenizer struct {
	Filename       string
	Options        TokenizerOptions
	tokenManager   *tokenManager
	keywordManager *keywordManager
}
