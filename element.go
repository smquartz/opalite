package anitogo

// elementCategory is a self-explanatory enum type
type elementCategory string

// elementCategory enumurations
const (
	elementCategoryAnimeSeason         elementCategory = "anime_season"
	elementCategoryAnimeSeasonPrefix   elementCategory = "anime_season_prefix"
	elementCategoryAnimeTitle          elementCategory = "anime_title"
	elementCategoryAnimeType           elementCategory = "anime_type"
	elementCategoryAnimeYear           elementCategory = "anime_year"
	elementCategoryAudioTerm           elementCategory = "audio_term"
	elementCategoryDeviceCompatibility elementCategory = "device_compatibility"
	elementCategoryEpisodeNumber       elementCategory = "episode_number"
	elementCategoryEpisodeNumberAlt    elementCategory = "episode_number_alt"
	elementCategoryEpisodePrefix       elementCategory = "episode_prefix"
	elementCategoryEpisodeTitle        elementCategory = "episode_title"
	elementCategoryFileChecksum        elementCategory = "file_checksum"
	elementCategoryFileExtension       elementCategory = "file_extension"
	elementCategoryFileName            elementCategory = "file_name"
	elementCategoryLanguage            elementCategory = "language"
	elementCategoryOther               elementCategory = "other"
	elementCategoryReleaseGroup        elementCategory = "release_group"
	elementCategoryReleaseInformation  elementCategory = "release_information"
	elementCategoryReleaseVersion      elementCategory = "release_version"
	elementCategorySource              elementCategory = "source"
	elementCategorySubtitles           elementCategory = "subtitles"
	elementCategoryVideoResolution     elementCategory = "video_resolution"
	elementCategoryVideoTerm           elementCategory = "video_term"
	elementCategoryVolumeNumber        elementCategory = "volume_number"
	elementCategoryVolumePrefix        elementCategory = "volume_prefix"
	elementCategoryUnknown             elementCategory = "unknown"
)

func elementCategoryInSlice(slice []elementCategory, cat elementCategory) bool {
	for _, v := range slice {
		if cat == v {
			return true
		}
	}
	return false
}

func (ec elementCategory) isSearchable() bool {
	searchableCategories := []elementCategory{
		elementCategoryAnimeSeasonPrefix,
		elementCategoryAnimeType,
		elementCategoryAudioTerm,
		elementCategoryDeviceCompatibility,
		elementCategoryEpisodePrefix,
		elementCategoryFileChecksum,
		elementCategoryLanguage,
		elementCategoryOther,
		elementCategoryReleaseGroup,
		elementCategoryReleaseInformation,
		elementCategoryReleaseVersion,
		elementCategorySource,
		elementCategorySubtitles,
		elementCategoryVideoResolution,
		elementCategoryVideoTerm,
		elementCategoryVolumePrefix,
	}
	return elementCategoryInSlice(searchableCategories, ec)
}

func (ec elementCategory) isSingular() bool {
	nonSingularCategories := []elementCategory{
		elementCategoryAnimeSeason,
		elementCategoryAnimeType,
		elementCategoryAudioTerm,
		elementCategoryDeviceCompatibility,
		elementCategoryEpisodeNumber,
		elementCategoryLanguage,
		elementCategoryOther,
		elementCategoryReleaseInformation,
		elementCategorySource,
		elementCategoryVideoTerm,
	}
	return !elementCategoryInSlice(nonSingularCategories, ec)
}

// this point onwards is likely misinterpreted

var elems *elements

type elements struct {
	elements       map[elementCategory]string
	checkAltNumber bool
}

func (e *elements) insert(cat elementCategory, content string) {
	e.elements[cat] = content
}

func (e elements) contains(cat elementCategory) bool {
	_, has := e.elements[cat]
	return has
}

func (e elements) empty() bool {
	return len(e.elements) <= 0
}
