package anitogo

import (
	"strings"

	"github.com/juju/errors"
)

type keywordOption struct {
	unidentifiable bool
	unsearchable   bool
	invalid        bool
}

type keyword struct {
	category elementCategory
	options  keywordOption
}

func (k keyword) empty() bool {
	var emptyKeyword keyword
	return k == emptyKeyword
}

type keywordManager struct {
	// optionsDefault
	// optionsInvalid
	// optionsUnidentifiable
	// optionsUnidentifiableInvalid
	// optionsUnidentifiableUnsearchable
	fileExtensions map[string]keyword
	keys           map[string]keyword
}

var keyman *keywordManager

func (km *keywordManager) addKeyword(cat elementCategory, options keywordOption, kwds []string) {
	for _, kw := range kwds {
		if _, alreadyIn := km.keys[kw]; alreadyIn || kw == "" {
			continue
		}

		if cat == elementCategoryFileExtension {
			km.fileExtensions[kw] = keyword{category: cat, options: options}
		} else {
			km.keys[kw] = keyword{category: cat, options: options}
		}
	}
}

type preidentifiedToken struct {
	beginIndex int
	endIndex   int
}

// wtf does this even do
func (km *keywordManager) peek(str string) []preidentifiedToken {
	/* this is a ported implementation
	type entry struct {
		category elementCategory
		keywords []string
	}
			entries := []entry{
				entry{category: elementCategoryAudioTerm, keywords: []string{"Dual Audio"}},
				entry{category: elementCategoryVideoTerm, keywords: []string{"H264", "H.264", "h264", "h.264"}},
				entry{category: elementCategoryVideoResolution, keywords: []string{"480p", "720p", "1080p"}},
				entry{category: elementCategorySource, keywords: []string{"Blu-Ray"}},
			}

		var preidentifiedTokens []preidentifiedToken

		for _, entry := range entries {
			for _, kw := range entry.keywords {
				keywordBeginIndex := strings.Index(str, kw)
				if keywordBeginIndex > -1 {
					elems.insert(entry.category, kw)

					keywordEndIndex := keywordBeginIndex + len(kw)
					preidentifiedTokens = append(preidentifiedTokens, preidentifiedToken{beginIndex: keywordBeginIndex, endIndex: keywordEndIndex})
				}
			}
		}
	*/

	// this i think is a better implementation
	var preidentifiedTokens []preidentifiedToken

	for kw, kwdetails := range keyman.keys {
		keywordBeginIndex := strings.Index(str, kw)
		if keywordBeginIndex > -1 {
			elems.insert(kwdetails.category, kw)

			keywordEndIndex := keywordBeginIndex + len(kw)
			preidentifiedTokens = append(preidentifiedTokens, preidentifiedToken{beginIndex: keywordBeginIndex, endIndex: keywordEndIndex})
		}
	}

	return preidentifiedTokens
}

func initKeywordManager() (*keywordManager, error) {
	km := new(keywordManager)
	km.fileExtensions = make(map[string]keyword)
	km.keys = make(map[string]keyword)

	optionsDefault := keywordOption{}
	optionsInvalid := keywordOption{invalid: true}
	optionsUnidentifiable := keywordOption{unidentifiable: true}
	optionsUnidentifiableInvalid := keywordOption{unidentifiable: true, invalid: true}
	optionsUnidentifiableUnsearchable := keywordOption{unidentifiable: true, unsearchable: true}

	km.addKeyword(elementCategoryAnimeSeasonPrefix, optionsUnidentifiable, []string{"SAISON", "SEASON"})
	km.addKeyword(elementCategoryAnimeType, optionsUnidentifiable, []string{"GEKIJOUBAN", "MOVIE",
		"OAD", "OAV", "ONA", "OVA",
		"SPECIAL", "SPECIALS",
		"TV",
	})
	km.addKeyword(elementCategoryAnimeType, optionsUnidentifiableUnsearchable, []string{"SP"})
	km.addKeyword(elementCategoryAnimeType, optionsUnidentifiableInvalid, []string{"ED", "ENDING", "NCED",
		"NCOP", "OP", "OPENING",
		"PREVIEW", "PV",
	})
	km.addKeyword(elementCategoryAudioTerm, optionsDefault, []string{
		// Audio channels
		"2.0CH", "2CH", "5.1", "5.1CH", "DTS", "DTS-ES", "DTS5.1",
		"TRUEHD5.1",
		// Audio codec
		"AAC", "AACX2", "AACX3", "AACX4", "AC3", "EAC3", "E-AC-3",
		"FLAC", "FLACX2", "FLACX3", "FLACX4", "LOSSLESS", "MP3", "OGG",
		"VORBIS",
		// Audio language
		"DUALAUDIO", "DUAL AUDIO",
	})
	km.addKeyword(elementCategoryDeviceCompatibility, optionsDefault, []string{"IPAD3", "IPHONE5", "IPOD", "PS3", "XBOX", "XBOX360"})
	km.addKeyword(elementCategoryDeviceCompatibility, optionsUnidentifiable, []string{"ANDROID"})
	km.addKeyword(elementCategoryEpisodePrefix, optionsDefault, []string{"EP", "EP.", "EPS", "EPS.", "EPISODE", "EPISODE.", "EPISODES",
		"CAPITULO", "EPISODIO", "FOLGE",
	})
	km.addKeyword(elementCategoryEpisodePrefix, optionsInvalid, []string{"E", "\x7B2C"})
	km.addKeyword(elementCategoryFileExtension, optionsDefault, []string{"3GP", "AVI", "DIVX", "FLV", "M2TS", "MKV", "MOV", "MP4", "MPG",
		"OGM", "RM", "RMVB", "WEBM", "WMV",
	})
	km.addKeyword(elementCategoryFileExtension, optionsInvalid, []string{"AAC", "AIFF", "FLAC", "M4A", "MP3", "MKA", "OGG", "WAV", "WMA",
		"7Z", "RAR", "ZIP",
		"ASS", "SRT",
	})
	km.addKeyword(elementCategoryLanguage, optionsDefault, []string{"ENG", "ENGLISH", "ESPANOL", "JAP", "PT-BR", "SPANISH", "VOSTFR"})
	km.addKeyword(elementCategoryLanguage, optionsUnidentifiable, []string{"ESP", "ITA"})
	km.addKeyword(elementCategoryOther, optionsDefault, []string{"REMASTER", "REMASTERED", "UNCENSORED", "UNCUT",
		"TS", "VFR", "WIDESCREEN", "WS",
	})
	km.addKeyword(elementCategoryReleaseGroup, optionsDefault, []string{"THORA"})
	km.addKeyword(elementCategoryReleaseInformation, optionsDefault, []string{"BATCH", "COMPLETE", "PATCH", "REMUX"})
	km.addKeyword(elementCategoryReleaseInformation, optionsUnidentifiable, []string{"END", "FINAL"})
	km.addKeyword(elementCategorySource, optionsDefault, []string{
		"BD", "BDRIP", "BLURAY", "BLU-RAY",
		"DVD", "DVD5", "DVD9", "DVD-R2J", "DVDRIP", "DVD-RIP",
		"R2DVD", "R2J", "R2JDVD", "R2JDVDRIP",
		"HDTV", "HDTVRIP", "TVRIP", "TV-RIP",
		"WEBCAST", "WEBRIP",
	})
	km.addKeyword(elementCategorySubtitles, optionsDefault, []string{
		"ASS", "BIG5", "DUB", "DUBBED", "HARDSUB", "HARDSUBS", "RAW",
		"SOFTSUB", "SOFTSUBS", "SUB", "SUBBED", "SUBTITLED",
	})
	km.addKeyword(elementCategoryVideoTerm, optionsDefault, []string{
		// Frame rate
		"23.976FPS", "24FPS", "29.97FPS", "30FPS", "60FPS", "120FPS",
		// Video codec
		"8BIT", "8-BIT", "10BIT", "10BITS", "10-BIT", "10-BITS",
		"HI10", "HI10P", "HI444", "HI444P", "HI444PP",
		"H264", "H265", "H.264", "H.265", "X264", "X265", "X.264",
		"AVC", "HEVC", "HEVC2", "DIVX", "DIVX5", "DIVX6", "XVID",
		// Video format
		"AVI", "RMVB", "WMV", "WMV3", "WMV9",
		// Video quality
		"HQ", "LQ",
		// Video resolution
		"HD", "SD",
	})
	km.addKeyword(elementCategoryVolumePrefix, optionsDefault, []string{"VOL", "VOL.", "VOLUME"})

	return km, errors.Trace(nil)
}

func init() {
	var err error
	keyman, err = initKeywordManager()
	if err != nil {
		panic(err)
	}
}
