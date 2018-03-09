package anitogo

func (k keywordDescriptor) Empty() bool {
	return k == keywordDescriptor{}
}

func (k keywordManager) initialised() bool {
	return k.keywords != nil && k.fileExtensions != nil
}
