package consumer

import (
	"bytes"
	"os"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/collections"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/turbine/queue"
	"github.com/rs/zerolog/log"
)

func create(downloadsFolder string, stream mapof.Any) queue.Result {

	const location = "consumer.create"

	// Create a new ActivityStreams client
	client := streams.NewDefaultClient()
	streamID := stream.GetString("streamId")

	// Load the Album from the provided URL
	album, err := client.Load(stream.GetString("url"))

	if err != nil {
		return queue.Error(derp.Wrap(err, location, "Unable to load album"))
	}

	// Load the Tracks for the Album
	tracks, err := client.Load(album.Get("tracks").String())

	if err != nil {
		return queue.Error(derp.Wrap(err, location, "Unable to load tracks for album", album.ID()))
	}

	// If the album is empty, then there's nothing more to do
	if tracks.TotalItems() == 0 {
		return queue.Success()
	}

	// For each track in the album, download the primary attachment file
	counter := 1
	for track := range collections.RangeDocuments(tracks) {

		// Only process audio tracks (why would anything else be in an album?)
		if track.Type() != vocab.ObjectTypeAudio {
			continue
		}

		// Get/validate the media document
		media := track.Get(vocab.PropertyURL)
		if media.Type() != vocab.CoreTypeLink {
			continue
		}

		// Href must be valid
		if media.Href() == "" {
			continue
		}

		// Download the audio file
		var buffer bytes.Buffer
		if err := remote.Get(media.Href()).Result(&buffer).Send(); err != nil {
			derp.Report((derp.Wrap(err, location, "Unable to download track", media.Href())))
			continue
		}

		albumFolder := downloadsFolder + "/" + streamID

		if counter == 1 {
			// Create a directory for the album tracks
			if err := os.Mkdir(albumFolder, os.ModePerm); err != nil {
				log.Debug().Str("directory", albumFolder).Msg("Unable to create directory (it probably already exists)")
			}
		}

		// Try to save the track to disk
		filename := albumFolder + "/" + trackFilename(track, counter)
		if err := os.WriteFile(filename, buffer.Bytes(), os.ModePerm); err != nil {
			derp.Report(derp.Wrap(err, location, "Unable to save track", filename))
			continue
		}

		// Increment track counter
		counter++
	}

	if counter == 1 {
		log.Info().Str("album", album.Name()).Int("tracks", counter-1).Msg("Album is empty")
	} else {
		log.Info().Str("album", album.Name()).Int("tracks", counter-1).Msg("Album imported")
	}

	return queue.Success()
}
