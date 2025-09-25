package consumer

import (
	"os"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/turbine/queue"
)

func delete(downloadsFolder string, stream mapof.Any) queue.Result {

	const location = "consumer.create"

	// Calculate the name of the folder to delete
	albumFolder := downloadsFolder + "/" + stream.GetString("streamId")

	// Try to delete the album folder and all of its contents
	if err := os.RemoveAll(albumFolder); err != nil {
		return queue.Error(derp.Wrap(err, location, "Unable to delete album folder", albumFolder))
	}

	return queue.Success()
}
