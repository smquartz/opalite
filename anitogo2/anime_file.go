package anitogo2

import "reflect"

// AnimeFile describes a file that we are being asked te parse the name of
type AnimeFile struct {
	AnimeSeason         []uint
	AnimeSeasonPrefix   string // S
	AnimeTitle          string
	AnimeType           []string // S
	AnimeYear           uint
	AudioTerm           []string // S
	DeviceCompatibility []string // S
	EpisodeNumber       []uint
	EpisodeNumberAlt    uint
	EpisodePrefix       string // S
	EpisodeTitle        string
	FileChecksum        string //S
	FileExtension       string
	FileName            string
	Language            []string // S
	Other               []string // S
	ReleaseGroup        string   // S
	ReleaseInformation  []string // S
	ReleaseVersion      uint     // S
	Source              []string // S
	Subtitles           string   // S
	VideoResolution     string   //S
	VideoTerm           []string // S
	VolumeNumber        uint
	VolumePrefix        string //S
	Unknown             string
	checkAltNumber      bool
}

// fieldSearchable returns whether a given field is considered searchable or not
func (AnimeFile) fieldSearchable(field interface{}) bool {
	searchableFields := map[string]bool{
		"AnimeSeasonPrefix":   true,
		"AnimeType":           true,
		"AudioTerm":           true,
		"DeviceCompatibility": true,
		"EpisodePrefix":       true,
		"FileChecksum":        true,
		"Language":            true,
		"Other":               true,
		"ReleaseGroup":        true,
		"ReleaseInformation":  true,
		"ReleaseVersion":      true,
		"Source":              true,
		"Subtitles":           true,
		"VideoResolution":     true,
		"VideoTerm":           true,
		"VolumePrefix":        true,
	}

	indirect := reflect.Indirect(reflect.ValueOf(field))
	fieldName := indirect.Type().Name()

	searchable, found := searchableFields[fieldName]
	return searchable && found
}
