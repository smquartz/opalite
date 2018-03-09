package anitogo

import "github.com/juju/errors"

func (k *keywordManager) SetKeyword(cat category, options keywordOptions, kwd keyword) error {
	if !k.initialised() {
		return errors.Trace(ErrorKeywordManagerMapsNotInitialized)
	}

	if cat == categoryFileExtension {
		k.fileExtensions[kwd] = keywordDescriptor{Category: cat, Options: options}
	} else {
		k.keywords[kwd] = keywordDescriptor{Category: cat, Options: options}
	}

	return errors.Trace(nil)
}

func (k *keywordManager) BulkSetKeywords(cat category, options keywordOptions, kwds []keyword) error {
	for i, kwd := range kwds {
		err := k.SetKeyword(cat, options, kwd)
		if err != nil {
			return errors.Trace(errors.Annotatef(err, "BulkSetKeywords keyword %v", i))
		}
	}
	return errors.Trace(nil)
}

func (k keywordManager) GetKeyword(kwd keyword) (keywordDescriptor, error) {
	if !k.initialised() {
		return keywordDescriptor{}, errors.Trace(ErrorKeywordManagerMapsNotInitialized)
	}

	kd, found := k.keywords[kwd]
	if !found {
		kd, found = k.fileExtensions[kwd]
		if !found {
			return keywordDescriptor{}, errors.Trace(ErrorKeywordNotFound)
		}
		return kd, errors.Trace(nil)
	}
	return kd, errors.Trace(nil)
}
