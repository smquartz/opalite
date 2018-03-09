package anitogo

// describes a category of token
type tokenCategory int

// describes a flags set on a token to describe it
type tokenFlags int

// tokenCategory enums
const (
	tokenCategoryUnknown tokenCategory = 1 << iota
	tokenCategoryBracket
	tokenCategoryDelimiter
	tokenCategoryIdentifier
	tokenCategoryInvalid
)

// tokenFlags enums
const (
	tokenFlagsNone    tokenFlags = 0
	tokenFlagsBracket tokenFlags = 1 << iota
	tokenFlagsNotBracket
	tokenFlagsDelimiter
	tokenFlagsNotDelimiter
	tokenFlagsIdentifier
	tokenFlagsNotIdentifier
	tokenFlagsUnknown
	tokenFlagsNotUnknown
	tokenFlagsValid
	tokenFlagsNotValid
	tokenFlagsEnclosed
	tokenFlagsNotEnclosed
)

// tokenFlags mask enums
const (
	tokenFlagsMaskCategories = tokenFlagsBracket | tokenFlagsNotBracket | tokenFlagsDelimiter | tokenFlagsNotDelimiter | tokenFlagsIdentifier | tokenFlagsNotIdentifier | tokenFlagsUnknown | tokenFlagsNotUnknown | tokenFlagsValid | tokenFlagsNotValid
	toeknFlagsMaskEnclosed   = tokenFlagsEnclosed | tokenFlagsNotEnclosed
)

// token describes a logical segment of a filename
type token struct {
	Category tokenCategory
	Flags    tokenFlags
	Content  string
	Enclosed bool
}

// tokens describes a slice of token pointers
type tokens []*token

// tokenManager describes a "manager" struct that stores token pointers
type tokenManager struct {
	tokens tokens
}

// indexSet describes a pair of start and finish indexes
type indexSet struct {
	Begin int
	End   int
}
