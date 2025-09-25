package consumer

import (
	"regexp"

	"github.com/benpate/hannibal/streams"
	"github.com/benpate/rosetta/convert"
)

var filenameRegexp = regexp.MustCompile("[^a-zA-Z0-9]+")

func trackFilename(track streams.Document, index int) string {

	trackName := track.Name()
	trackName = filenameRegexp.ReplaceAllString(trackName, "_")
	trackName = convert.String(index) + "_" + trackName + ".mp3"

	return trackName
}
