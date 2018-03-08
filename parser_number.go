package anitogo

import (
	"strconv"

	"github.com/juju/errors"
)

// Some validation constants
const (
	animeYearMin     = 1900
	animeYearMax     = 2050
	episodeNumberMax = animeYearMin - 1
	volumeNumberMax  = 20 // is this too low?
)

func validEpisodeNumber(numstr string) bool {
	num, err := strconv.ParseUint(numstr, 10, 16)
	if err != nil {
		return false
	}
	return num <= episodeNumberMax
}

func setEpisodeNumber(numstr string, tkn token, validate bool) error {
	if validate && !validEpisodeNumber(numstr) {
		return errors.Trace(ErrorInvalidEpisodeNumber)
	}

	tkn.category = tokenCategoryIdentifier

	category := elementCategoryEpisodeNumber

	// handle equivalent numbers
	if elems.checkAltNumber {
		// TODO: check if getting only the first episode number is enough
		episodeNumberStr := elems.elements[elementCategoryEpisodeNumber][0]
		num, err := strconv.ParseUint(numstr, 10, 16)
		if err != nil {
			return errors.Trace(ErrorInvalidEpisodeNumber)
		}
		episodeNumber, err := strconv.ParseUint(episodeNumberStr, 10, 16)
		if err != nil {
			return errors.Trace(ErrorInvalidEpisodeNumber)
		}

		if num > episodeNumber {
			category = elementCategoryEpisodeNumberAlt
		} else if num < episodeNumber {
			delete(elems.elements, elementCategoryEpisodeNumber)
			elems.insert(elementCategoryEpisodeNumberAlt, episodeNumberStr)
		}
	}

	elems.insert(category, numstr)
	return errors.Trace(nil)
}

func setAlternativeEpisodeNumber(numstr string, tkn *token) error {
	elems.insert(elementCategoryEpisodeNumberAlt, numstr)
	tkn.category = tokenCategoryIdentifier

	return errors.Trace(nil)
}

func (t tokens) checkExtentKeyword(cat elementCategory, tkn *token) error {
	nextToken, err := t.findNext(*tkn, tokenFlagNotDelimiter)
	if err != nil {
		return errors.Trace(err)
	}

	if nextToken.category == tokenCategoryUnknown {
		if !nextToken.empty() && findNumberInString(nextToken.content) > -1 {
			if cat == elementCategoryEpisodeNumber {
				if !match
			}
		}
	}

	return errors.Trace(nil)
}
