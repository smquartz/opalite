package anitogo

func newKeywordManager() *keywordManager {
	km := new(keywordManager)
	km.keywords = make(keywords)
	km.fileExtensions = make(keywords)

	optionsDefault := keywordOptions{}
	optionsInvalid := keywordOptions{Invalid: true}
	optionsUnidentifiable := keywordOptions{Unidentifiable: true}
	optionsUnidentifiableInvalid := keywordOptions{Unidentifiable: true, Invalid: true}
	optionsUnidentifiableUnsearchable := keywordOptions{Unidentifiable: true, Unsearchable: true}

	km.BulkSetKeywords(categoryAnimeSeasonPrefix, optionsUnidentifiable, []keyword{"SAISON", "SEASON"})
	km.BulkSetKeywords(categoryAnimeType, optionsUnidentifiable, []keyword{"GEKIJOUBAN", "MOVIE",
		"OAD", "OAV", "ONA", "OVA",
		"SPECIAL", "SPECIALS",
		"TV",
	})
	km.BulkSetKeywords(categoryAnimeType, optionsUnidentifiableUnsearchable, []keyword{"SP"})
	km.BulkSetKeywords(categoryAnimeType, optionsUnidentifiableInvalid, []keyword{"ED", "ENDING", "NCED",
		"NCOP", "OP", "OPENING",
		"PREVIEW", "PV",
	})
	km.BulkSetKeywords(categoryAudioTerm, optionsDefault, []keyword{
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
	km.BulkSetKeywords(categoryDeviceCompatibility, optionsDefault, []keyword{"IPAD3", "IPHONE5", "IPOD", "PS3", "XBOX", "XBOX360"})
	km.BulkSetKeywords(categoryDeviceCompatibility, optionsUnidentifiable, []keyword{"ANDROID"})
	km.BulkSetKeywords(categoryEpisodePrefix, optionsDefault, []keyword{"EP", "EP.", "EPS", "EPS.", "EPISODE", "EPISODE.", "EPISODES",
		"CAPITULO", "EPISODIO", "FOLGE",
	})
	km.BulkSetKeywords(categoryEpisodePrefix, optionsInvalid, []keyword{"E", "\x7B2C"})
	km.BulkSetKeywords(categoryFileExtension, optionsDefault, []keyword{"3GP", "AVI", "DIVX", "FLV", "M2TS", "MKV", "MOV", "MP4", "MPG",
		"OGM", "RM", "RMVB", "WEBM", "WMV",
	})
	km.BulkSetKeywords(categoryFileExtension, optionsInvalid, []keyword{"AAC", "AIFF", "FLAC", "M4A", "MP3", "MKA", "OGG", "WAV", "WMA",
		"7Z", "RAR", "ZIP",
		"ASS", "SRT",
	})
	km.BulkSetKeywords(categoryLanguage, optionsDefault, []keyword{"ENG", "ENGLISH", "ESPANOL", "JAP", "PT-BR", "SPANISH", "VOSTFR"})
	km.BulkSetKeywords(categoryLanguage, optionsUnidentifiable, []keyword{"ESP", "ITA"})
	km.BulkSetKeywords(categoryOther, optionsDefault, []keyword{"REMASTER", "REMASTERED", "UNCENSORED", "UNCUT",
		"TS", "VFR", "WIDESCREEN", "WS",
	})
	km.BulkSetKeywords(categoryReleaseGroup, optionsDefault, []keyword{"THORA"})
	km.BulkSetKeywords(categoryReleaseInformation, optionsDefault, []keyword{"BATCH", "COMPLETE", "PATCH", "REMUX"})
	km.BulkSetKeywords(categoryReleaseInformation, optionsUnidentifiable, []keyword{"END", "FINAL"})
	km.BulkSetKeywords(categorySource, optionsDefault, []keyword{
		"BD", "BDRIP", "BLURAY", "BLU-RAY",
		"DVD", "DVD5", "DVD9", "DVD-R2J", "DVDRIP", "DVD-RIP",
		"R2DVD", "R2J", "R2JDVD", "R2JDVDRIP",
		"HDTV", "HDTVRIP", "TVRIP", "TV-RIP",
		"WEBCAST", "WEBRIP",
	})
	km.BulkSetKeywords(categorySubtitles, optionsDefault, []keyword{
		"ASS", "BIG5", "DUB", "DUBBED", "HARDSUB", "HARDSUBS", "RAW",
		"SOFTSUB", "SOFTSUBS", "SUB", "SUBBED", "SUBTITLED",
	})
	km.BulkSetKeywords(categoryVideoTerm, optionsDefault, []keyword{
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
	km.BulkSetKeywords(categoryVolumePrefix, optionsDefault, []keyword{"VOL", "VOL.", "VOLUME"})

	return km
}
