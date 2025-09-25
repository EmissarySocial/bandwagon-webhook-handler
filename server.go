package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/EmissarySocial/bandwagon-webhook-handler/config"
	"github.com/EmissarySocial/bandwagon-webhook-handler/consumer"
	"github.com/EmissarySocial/bandwagon-webhook-handler/handler"
	"github.com/benpate/derp"
	"github.com/benpate/digital-dome/dome"
	"github.com/benpate/digital-dome/dome4echo"
	"github.com/benpate/domain"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/turbine/queue"
	"github.com/benpate/turbine/queue_filesystem"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func main() {

	log.Info().Msg("Starting Bandwagon Webhook Handler")

	// Parse command line arguments
	args := config.GetCommandLineArgs()

	// Create the task queue to process inbound webhook events
	q := queue.New(
		queue.WithConsumers(consumer.New(args.Downloads)),   // Consumer executes tasks
		queue.WithStorage(queue_filesystem.New(args.Queue)), // Storage saves tasks
		queue.WithWorkerCount(args.Workers),                 // Worker count should be small
		queue.WithBufferSize(args.Workers),                  // Match buffer size to worker count
	)

	// Configure the web server
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.HTTPErrorHandler = errorHandler

	// Configure digital dome (Application Firewall)
	digitalDome := dome.New()
	e.Use(dome4echo.New(&digitalDome))

	// Configure Server Routes
	e.GET("/", handler.GetPage())
	e.POST("/", handler.PostPage(args, q))

	// Start HTTP servers
	if args.HTTPPort > 0 {
		go derp.Report(e.Start(":" + convert.String(args.HTTPPort)))
	}

	// Start HTTPS server (if configured)
	if args.HTTPSPort > 0 {
		go derp.Report(e.StartAutoTLS(":" + convert.String(args.HTTPSPort)))
	}

	// Listen to the OS SIGINT channel for an interrupt signal
	// Use a buffered channel to avoid missing signals as recommended
	// https://golang.org/pkg/os/signal/#Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Wait for the "quit" signal from the OS, then shut down
	<-quit

	// Get a cancellation context with a 30 second timeout, which
	// hopefully lets us finish processing any in-flight requests
	// https://echo.labstack.com/docs/cookbook/graceful-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Try to shut down the server
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	// Last, stop the task queue (which waits for all workers to finish their tasks)
	q.Stop()
}

// errorHandler is a custom error handler that returns a JSON error message to the client
func errorHandler(err error, ctx echo.Context) {

	// Write the error to the console (on production and local domains)
	derp.Report(err)

	// On localhost, allow developers to see full error dump.
	if domain.IsLocalhost(domain.Hostname(ctx.Request())) {
		_ = ctx.JSONPretty(derp.ErrorCode(err), err, "  ")
		return
	}

	// Fall through to general error handler
	_ = ctx.String(derp.ErrorCode(err), derp.Message(err))
}
