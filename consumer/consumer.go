package consumer

import (
	"github.com/benpate/derp"
	"github.com/benpate/turbine/queue"
)

func New(directory string) queue.Consumer {

	const location = "consumer.Consumer"

	return func(name string, args map[string]any) queue.Result {

		switch name {

		case "Create":
			return create(directory, args)

		case "Update":
			return create(directory, args)

		case "Delete":
			return delete(directory, args)
		}

		return queue.Failure(derp.NotFound(location, "Unrecognized task name", name))
	}
}
