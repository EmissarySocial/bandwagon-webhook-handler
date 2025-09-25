// Package config reads command line arguments and returns a configuration object
// to the caller
package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
)

// CommandLineArgs represents the command line arguments passed to the server
type CommandLineArgs struct {
	Downloads string // Location to store downloaded files
	Queue     string // Locate to store task queue
	Workers   int    // Number of worker threads to process the queue
	HTTPPort  int    // Port for HTTP server
	HTTPSPort int    // Port for HTTPS server
}

// GetCommandLineArgs returns the location of the configuration file
func GetCommandLineArgs() CommandLineArgs {

	// Default configuration
	result := CommandLineArgs{}

	// Look for the configuration location in the command line arguments
	pflag.StringVar(&result.Downloads, "export", "./downloads", "Path to save downloaded files. (default => ./downloads)")
	pflag.StringVar(&result.Queue, "queue", "./queue", "Path to save task queue. (default => ./queue)")
	pflag.IntVar(&result.Workers, "workers", 1, "Number of worker threads to process the queue. (default => 1)")
	pflag.IntVar(&result.HTTPPort, "port", 8080, "Port for HTTP server. (default => 8080)")
	pflag.IntVar(&result.HTTPSPort, "https", 0, "Port for HTTPS server. (default => 443)")
	pflag.Parse()

	log.Debug().Msg("Reading command line arguments")
	log.Debug().Str("export", result.Downloads).Msg("export folder")
	log.Debug().Str("queue", result.Queue).Msg("queue folder")
	log.Debug().Int("workers", result.Workers).Msg("worker threads")
	log.Debug().Int("httpPort", result.HTTPPort).Msg("HTTP port")
	log.Debug().Int("httpsPort", result.HTTPSPort).Msg("HTTPS port")

	// Success!
	return result
}
