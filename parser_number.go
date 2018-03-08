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
	if elems.checkAltNumber() {
		// TODO: check if getting only the first episode number is enough
		episodeNumber := elems.elements
	}
}
