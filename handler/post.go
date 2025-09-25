package handler

import (
	"github.com/EmissarySocial/bandwagon-webhook-handler/config"
	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/turbine/queue"
	"github.com/labstack/echo/v4"
)

func PostPage(args config.CommandLineArgs, q *queue.Queue) echo.HandlerFunc {

	const location = "handler.PostPage"

	return func(ctx echo.Context) error {

		// Collect the data from the document body
		data := mapof.NewAny()
		if err := ctx.Bind(&data); err != nil {
			return derp.Wrap(err, location, "Unable to parse webhook payload")
		}

		// Pass the object to the queue based on the event type
		switch data.GetString("type") {

		case "Create":
			q.Enqueue <- queue.NewTask("Create", data.GetMap("object"))
			return ctx.NoContent(200)

		case "Update":
			q.Enqueue <- queue.NewTask("Update", data.GetMap("object"))
			return ctx.NoContent(200)

		case "Delete":
			q.Enqueue <- queue.NewTask("Delete", data.GetMap("object"))
			return ctx.NoContent(200)

		}

		// Unrecognized event type
		return derp.BadRequest(location, "Event type must be 'Create', 'Update', or 'Delete'")
	}
}
