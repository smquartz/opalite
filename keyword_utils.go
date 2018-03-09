package anitogo

func (k keywordDescriptor) Empty() bool {
	return k == keywordDescriptor{}
}

func (k keywordManager) Initialised() bool {
	return k.keywords != nil && k.fileExtensions != nil
}
