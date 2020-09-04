package cmd

import (
	"log"

	"github.com/j4ng5y/dhwh/pkg/server"
	"github.com/spf13/cobra"
)

// Run is the primary function that runs the command line interface.
//
// Arguments:
//     None
//
// Returns:
//     None
func Run() {
	var (
		configFile string
		rootCMD    = &cobra.Command{
			Use:     "dhwh",
			Version: "0.0.1",
			Short:   "a simple webhook server for listening for, and replying to, dockerhub event pushes.",
			Run: func(ccmd *cobra.Command, args []string) {
				var S *server.Server
				if configFile == "" {
					S = server.NewWithOptions(server.WithDefaults())
				} else {
					C := new(Config)
					if err := C.Unmarshal(configFile); err != nil {
						log.Fatal(err)
					}

					rhto, err := ParseDuration(C.HTTPServer.ReadHeaderTimeout)
					if err != nil {
						log.Fatal(err)
					}

					rto, err := ParseDuration(C.HTTPServer.ReadTimeout)
					if err != nil {
						log.Fatal(err)
					}

					wto, err := ParseDuration(C.HTTPServer.WriteTimeout)
					if err != nil {
						log.Fatal(err)
					}

					ito, err := ParseDuration(C.HTTPServer.IdleTimeout)
					if err != nil {
						log.Fatal(err)
					}

					S = server.NewWithOptions(
						server.WithHTTPServerAddress(C.HTTPServer.IPAddress, C.HTTPServer.Port),
						server.WithHTTPServerReadHeaderTimeout(rhto),
						server.WithHTTPServerReadTimeout(rto),
						server.WithHTTPServerWriteTimeout(wto),
						server.WithHTTPServerIdleTimeout(ito))
				}

				S.Run()
			},
		}
	)

	rootCMD.Flags().StringVarP(&configFile, "config-file", "f", "", "The configuration file to use to run the server.")

	if err := rootCMD.Execute(); err != nil {
		log.Fatal(err)
	}
}
