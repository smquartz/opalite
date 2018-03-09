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

func (p *parser) SetEpisodeNumber(numstr string, tkn *token, validate bool) error {
	if validate && !validEpisodeNumber(numstr) {
		return errors.Trace(ErrorInvalidEpisodeNumber)
	}

	tkn.Category = tokenCategoryIdentifier

	category := categoryEpisodeNumber
	var episodeNumber uint
	var num uint

	// handle equivalent numbers
	if p.AnimeFile.checkAltNumber {
		// TODO: check if getting only the first episode number is enough
		episodeNumber = p.AnimeFile.EpisodeNumber[0]
		num64, err := strconv.ParseUint(numstr, 10, 16)
		if err != nil {
			return errors.Trace(ErrorInvalidEpisodeNumber)
		}
		num = uint(num64)

		if num > episodeNumber {
			category = categoryEpisodeNumberAlt
		} else if num < episodeNumber {
			p.AnimeFile.EpisodeNumber = []uint{}
			p.AnimeFile.EpisodeNumberAlt = episodeNumber
		}
	}

	if category == categoryEpisodeNumber {
		p.AnimeFile.EpisodeNumber = []uint{episodeNumber}
	} else if category == categoryEpisodeNumberAlt {
		p.AnimeFile.EpisodeNumberAlt = episodeNumber
	}

	return errors.Trace(nil)
}

func (p *parser) SetAlternativeEpisodeNumber(numstr string, tkn *token) error {
	num, err := strconv.ParseUint(numstr, 10, 16)
	if err != nil {
		return errors.Trace(err)
	}
	p.AnimeFile.EpisodeNumberAlt = uint(num)
	tkn.Category = tokenCategoryIdentifier

	return errors.Trace(nil)
}

func (p *parser) CheckExtentKeyword(cat category, tkn *token) error {
	nextToken, err := p.tokenizer.tokenManager.tokens.FindNext(*tkn, tokenFlagsNotDelimiter)
	if err != nil {
		return errors.Trace(err)
	}

	if nextToken.Category == tokenCategoryUnknown {
		if !nextToken.Empty() && findNumberInString(nextToken.Content) > -1 {
			if cat == categoryEpisodeNumber {
				if !p.MatchEpisodePattern(nextToken.Content, nextToken) {
					p.SetEpisodeNumber(nextToken.Content, nextToken, false)
				}
			} else if cat == categoryVolumeNumber {
				if !p.matchVolumePattern(nextToken.Content, nextToken) {
					p.SetVolumeNumber(nextToken.Content, nextToken, false)
				}
			} else {
				// not implemented?
				return errors.Trace(ErrorPlaceholder)
			}
			tkn.Category = tokenCategoryIdentifier
			return errors.Trace(nil)
		}
	}

	return errors.Trace(ErrorPlaceholder)
}
