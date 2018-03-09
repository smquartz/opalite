package anitogo

// keywordOptinos describes parameters of a keyword
type keywordOptions struct {
	Unidentifiable bool
	Unsearchable   bool
	Invalid        bool
}

// keywordDescriptor describes the category and parameters of a keyword
type keywordDescriptor struct {
	Category category
	Options  keywordOptions
}

// keyword describes a logical unit of the release name we are asked to parse
type keyword string

// keywords describes a map with keywords as keys, and their descriptors as values
type keywords map[keyword]keywordDescriptor

// keywordManager instances store the keywords for a particular parsing session
type keywordManager struct {
	keywords       keywords
	fileExtensions keywords
}
